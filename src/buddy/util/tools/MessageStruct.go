package tools

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	MSG_PACK_USE_JSON = true
	MSG_HEAD_LENGTH   = 4 + 4
)

var BinOrder = binary.LittleEndian

type Message struct {
	Version  int32
	Cmd      int32
	Data     []byte
	SendData interface{}
}

func NewMessage() *Message {
	this := &Message{}
	return this
}

// 将一个完整的消息包解析到本对象中
func (this *Message) Decode(data []byte) error {
	if data == nil || len(data) < MSG_HEAD_LENGTH {
		fmt.Errorf("Message.Decode() faild len = %d data = %v ，", len(data), data)
	}

	var err error
	length := len(data)
	buf := bytes.NewReader(data)

	// 读取消息头
	err = binary.Read(buf, BinOrder, &this.Version)
	if err != nil {
		return err
	}

	err = binary.Read(buf, BinOrder, &this.Cmd)
	if err != nil {
		return err
	}

	lessLen := length - MSG_HEAD_LENGTH
	if lessLen > 0 {
		this.Data = make([]byte, lessLen)
		copy(this.Data, data[MSG_HEAD_LENGTH:length])
	}
	return nil
}

// 将本对象压缩成bin流
func (this *Message) Encode() ([]byte, error) {
	var err error
	wBuf := bytes.NewBuffer(nil)

	// 写入消息头
	err = binary.Write(wBuf, BinOrder, this.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Write(wBuf, BinOrder, this.Cmd)
	if err != nil {
		return nil, err
	}
	// 写入Data
	if this.SendData != nil {
		var b []byte
		b, err = this.EncodeData(this.SendData)
		if err != nil {
			return nil, err
		}

		_, err = wBuf.Write(b)
		if err != nil {
			return nil, err
		}
	}
	// 返回二进制数据
	return wBuf.Bytes(), nil
}

func (this *Message) EncodeData(data interface{}) (b []byte, err error) {
	if data == nil {
		return nil, fmt.Errorf("Message.EncodeData() : data is nil")
	}
	if MSG_PACK_USE_JSON {
		b, err = json.Marshal(data)
	} else {
		//b, err = msgpack.Marshal(data)
	}
	return b, err
}

func (this *Message) DecodeData(data []byte, v interface{}) error {
	if data == nil {
		return fmt.Errorf("Message.DecodeData() : data is nil")
	}

	if v == nil {
		return fmt.Errorf("Message.DecodeData() : v is nil")
	}

	if MSG_PACK_USE_JSON {
		return json.Unmarshal(data, v)
	} else {
		//return msgpack.Unmarshal(data, v)
	}
	return nil
}

func (this *Message) PrintInfo() {
	fmt.Println("msg cmd:", strconv.FormatInt(int64(this.Cmd), 16))
	fmt.Println("msg ver:", this.Version)
	if this.Data != nil {
		type Animal struct {
			Id string
			Lv string
		}
		var player Animal
		iiserr := json.Unmarshal(this.Data, &player)
		if iiserr != nil {
			fmt.Println("json err is:", iiserr)
		} else {
			fmt.Printf("%+v\n", player)
		}
		fmt.Println("msg dat:", string(this.Data))
	}
}

func (this *Message) String() string {
	return fmt.Sprintf("Message(%+x)", this.Cmd)
}

func MessageError(msg *Message, text string, a ...interface{}) error {
	text = fmt.Sprintf(text, a...)
	if msg != nil {
		return fmt.Errorf("msg error:cmd=%d,version=%d,data=%v,errInfo=%s", msg.Cmd, msg.Version, string(msg.Data), text)
	}
	return fmt.Errorf(text)
}
