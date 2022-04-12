#include <iostream>
#include <queue>
#include <vector>
#include <algorithm>
#include "include/boost/thread/mutex.hpp"
 
template <typename T>
class ObjectPool
{
public:
    ObjectPool(size_t chunk_size = kdefault_size);
    ~ObjectPool();
 
    T* acquire_object();
    void release_object(T* obj);
 
protected:
    void allocate_chunk();
    inline static void destroy(T* obj);
    inline static T* construct();
 
private:
    std::queue<T*> _free_list;
    std::vector<T*> _all_objects;
    size_t chunk_size_;                         //对象池中预分配对象个数
    static const size_t kdefault_size = 25;     //默认对象池大小
    boost::mutex _mutex;                        //锁
 
    ObjectPool(const ObjectPool<T>& src);
    ObjectPool<T>& operator=(const ObjectPool<T>& rhs);
};


template <typename T>
T* ObjectPool<T>::construct()
{
    return new T;
}

template <typename T>
void ObjectPool<T>::destroy(T* obj)
{
    delete obj;
}

template <typename T>
ObjectPool<T>::ObjectPool(size_t chunk_size) : chunk_size_ (chunk_size)
{
    if (chunk_size_ <= 0)
    {
        std::cout << "Object size invalid" << std::endl;
    }
    else
    {
        for (size_t i = 0; i != chunk_size; ++i)
        {
            allocate_chunk();
        }
    }
}
 
template <typename T>
ObjectPool<T>::~ObjectPool()
{
    std::for_each(_all_objects.begin(), _all_objects.end(), destroy);
}
 
template <typename T>
void ObjectPool<T>::allocate_chunk()
{
    T* new_object = construct();
    _all_objects.push_back(new_object);
    _free_list.push(new_object);
}
 
template <typename T>
T* ObjectPool<T>::acquire_object()
{
    boost::mutex::scoped_lock scoped_lock(_mutex);
    if (_free_list.empty())
    {
        allocate_chunk();
    }
    T *obj = _free_list.front();
    _free_list.pop();
    return obj;
}
 
template <typename T>
void ObjectPool<T>::release_object(T* obj)
{
    boost::mutex::scoped_lock scoped_lock(_mutex);
    _free_list.push(obj);
}