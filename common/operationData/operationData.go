package operationData


import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"project/chatRoom/common/message"
)
//定义该类
type PackageOperation struct {
	//连接
	Conn net.Conn
	//缓冲
	Buf [8096]byte
}

func (this *PackageOperation)ReadPackage() (mes message.Message,err error){


	//conn.read 在conn没有被关闭的情况下才会被阻塞，先读取字节长度
	_,err = this.Conn.Read(this.Buf[:4])
	if err != nil{

		return
	}
	var packageLen uint32
	//转成uint32类型
	packageLen = binary.BigEndian.Uint32(this.Buf[:4])
	fmt.Println("接收的数据长度为",packageLen)

	//读取数据本身
	n,err := this.Conn.Read(this.Buf[:packageLen])
	fmt.Println("接收的数据内容为",string(this.Buf[:packageLen]))
	if n!= int(packageLen) || err != nil{
		return
	}
	//到此已经能接收到数据，反序列化数据
	err = json.Unmarshal(this.Buf[:packageLen],&mes)
	if err != nil {
		fmt.Printf("json.unmarsha err=%v",err )
		return
	}
	fmt.Println("读取数据成功")
	return
}
func (this *PackageOperation)WritePackage(data []byte) (err error)  {
	//先发送数据长度
	var packageLen uint32
	packageLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4],packageLen)
	//发送长度
	n,err :=this.Conn.Write(buf[:4])
	fmt.Println("发送的数据长度为",string(packageLen))
	if n!=4 || err!=nil {
		fmt.Println("conn.write dataLen fail err=",err)
		return
	}
	//发送数据本身
	n,err = this.Conn.Write(data)
	fmt.Println("发送的数据内容为",string(data))
	if n!= int(packageLen) || err != nil{
		fmt.Println("conn.write data fail err=",err)
		return
	}
	return
}
