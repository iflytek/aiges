package daemon

//
//type SegIdManager struct {
//	begin int64
//	end   int64
//	m     map[int64]bool //存放不能用的id
//	mu    sync.Mutex     //此处的锁用来保证返回的sedId不重复
//}
//
//func (s *SegIdManager) getMin() (min int64) {
//	s.mu.Lock()
//	defer func() {
//		s.m[min] = true
//		s.mu.Unlock()
//	}()
//
//	for startIx := s.begin; startIx <= s.end; startIx++ {
//		if _, ok := s.m[startIx]; !ok {
//			min = startIx
//			return
//		}
//	}
//	return atomic.AddInt64(&s.end, 1)
//}
//func (s *SegIdManager) free(in int64) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	delete(s.m, in)
//}
//
//func (s *SegIdManager) freeze(in int64) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	s.m[in] = true
//}
//
//var segIdManagerInst *SegIdManager
//
//func init() {
//	segIdManagerInst = &SegIdManager{begin: 0, m: make(map[int64]bool)}
//}
