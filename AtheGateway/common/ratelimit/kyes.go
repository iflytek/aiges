package ratelimit

type RateLimit interface {
	GetKey()string
	Climit()int
}


type AppidLimit struct {
	ConcurrencyLimit int `json:"limit"`
	Appid         string `json:"appid"`
	Policy        string `json:"policy"`
	FaultTolerant bool `json:"tolerant"`
	Enabled bool `json:"enabled"`
}



func (r *AppidLimit)Init() error {
	if r.ConcurrencyLimit <=0{
		r.ConcurrencyLimit = MaxData
	}
	return nil
}

func (r *AppidLimit) GetKey() string {

	return generateKey(r.Appid)
}

func (r *AppidLimit) Climit() int {
	return r.ConcurrencyLimit
}

type IpLimit struct {
	ConcurrencyLimit int `json:"limit"`
	Ip         string `json:"ip"`
	Policy        string `json:"policy"`
	FaultTolerant bool `json:"tolerant"`
	Enabled bool `json:"enabled"`
}

func (r *IpLimit) GetKey() string {
	return KeyPrefix+r.Ip
}

func (r *IpLimit)Climit()int  {
	return r.ConcurrencyLimit
}

type AppidConnLimit struct {
	ConcurrencyLimit int `json:"limit"`
	Appid         string `json:"appid"`
	Policy        string `json:"policy"`
	FaultTolerant bool `json:"tolerant"`
	Enabled bool `json:"enabled"`
}

func (r *AppidConnLimit)GetKey()  string{
	return generateConnKey(r.Appid)
}

func (r *AppidConnLimit)Climit()int  {
	return r.ConcurrencyLimit
}

type MaxConnLimit struct {
	ConcurrencyLimit int `json:"limit"`
	Policy        string `json:"policy"`
	FaultTolerant bool `json:"tolerant"`
	Enabled bool `json:"enabled"`
}

func (r *MaxConnLimit)GetKey()  string{
	return KeyConnMax
}

func (r *MaxConnLimit)Climit()int  {
	return r.ConcurrencyLimit
}