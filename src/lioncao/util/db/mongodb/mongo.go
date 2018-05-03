package mongodb

import (
	//"lioncao/util/tools"
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

const (
	RESULT_NOT_FOUND = "not found" // mongodb中未查取到的结果
)

type MgoCon struct {
	ip       string
	port     string
	username string
	password string
	session  *mgo.Session
	ok       bool // 是否可以进行数据操作了
}

func (this *MgoCon) Init(ip, port, username, password string) error {
	this.ip = ip
	this.port = port
	this.username = username
	this.password = password

	address := fmt.Sprintf("%s:%s", ip, port)
	var err error
	this.session, err = mgo.Dial(address)
	if err != nil {
		return err
	} else {
		this.session.SetMode(mgo.Monotonic, true)
	}

	if username != "" {
		err = this.Login(username, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *MgoCon) Close() {
	this.session.Close()
}

func (this *MgoCon) Insert(DB, collect string, data interface{}) error {
	c := this.session.DB(DB).C(collect)
	err := c.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) Update(DB, collect string, selector interface{}, data interface{}) error {
	c := this.session.DB(DB).C(collect)
	err := c.Update(selector, data)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) Find(DB, collect string, selector interface{}, result interface{}) error {
	c := this.session.DB(DB).C(collect)
	err := c.Find(selector).All(result)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) FindN(DB, collect string, N int, selector interface{}, result interface{}) error {
	c := this.session.DB(DB).C(collect)
	err := c.Find(selector).Limit(N).All(result)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) DropCollection(DB, collect string) error {
	c := this.session.DB(DB).C(collect)
	//err := c.Find(selector).All(result)
	err := c.DropCollection()
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) FindOne(DB, collect string, selector interface{}, result interface{}) error {
	c := this.session.DB(DB).C(collect)
	err := c.Find(selector).One(result)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) Remove(DB, collect string, selector interface{}) error {
	c := this.session.DB(DB).C(collect)
	err := c.Remove(selector)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoCon) Login(user, pass string) error {

	cred := &mgo.Credential{
		Username: user,
		Password: pass,
	}

	err := this.session.Login(cred)
	if err != nil {
		return err
	}
	return nil
}
