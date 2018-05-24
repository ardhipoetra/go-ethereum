package extension

import (
	"github.com/ethereum/go-ethereum/logger"
	"github.com/ethereum/go-ethereum/logger/glog"
	"os/exec"
	"time"
	"strconv"
	"os"
	"log"
	"bufio"
	"fmt"
)

type Interceptor struct {
	sendNode int
	recvNode int
	fileDir string
	filename string
	messageType string
	state int
	info int
	eventId int

	Exist chan bool
}
/*important global variable to decide whether intercept or not in the ethereum system*/
var IsIntercept bool = true

func TestIntercept(selfId int, messageType string, state int, info int, ipcdir string) *Interceptor{
	interceptor := &Interceptor{
		sendNode:		selfId,
		messageType:    messageType,
		state:			state,
		info:			info}

	interceptor.eventId = interceptor.gethash()
	interceptor.fileDir = ipcdir
	interceptor.filename = interceptor.createFilename()
	interceptor.Exist = make(chan bool)

	glog.V(logger.Info).Infoln("@RD state: ",state, "info: ", info)

	//start intercepting messages into the file inside /new/ folder
	newFileName := interceptor.fileDir+"/new/"+interceptor.filename
	newFile, err := os.Create(newFileName) // Truncates if file already exists, be careful!

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	writer := bufio.NewWriter(newFile)
	defer newFile.Close()

	recvNode := 0

	if selfId == 0 {
		recvNode = 1
	}

	fmt.Fprintln(writer, "sendNode=",selfId)
	fmt.Fprintln(writer, "recvNode=",recvNode) //dummy
	fmt.Fprintln(writer,"messageType=", messageType)
	fmt.Fprintln(writer, "pendingNonce=",state)
	fmt.Fprintln(writer, "currentNonce=",info)
	fmt.Fprintln(writer, "eventId=",interceptor.eventId)
	writer.Flush()

	interceptor.commitEvent()
	return interceptor
}

func NewIntercept(sendNode int, recvNode int, messageType string, state int, info int) *Interceptor{
	interceptor := &Interceptor{
		sendNode:		sendNode,
		recvNode:		recvNode,
		messageType:    messageType,
		state:			state,
		info:			info}

	interceptor.eventId = interceptor.gethash()
	interceptor.fileDir = "/tmp/ipc"
	interceptor.filename = interceptor.createFilename()
	interceptor.Exist = make(chan bool)

	glog.V(logger.Info).Infoln("@huanke state: ",state, "info: ", info)
	//start intercepting messages into the file inside /new/ folder
	newFileName := interceptor.fileDir+"/new/"+interceptor.filename
	newFile, err := os.Create(newFileName) // Truncates if file already exists, be careful!
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	writer := bufio.NewWriter(newFile)
	defer newFile.Close()
	fmt.Fprintln(writer, "sendNode=",sendNode)
	fmt.Fprintln(writer,"recvNode=",recvNode)
	fmt.Fprintln(writer,"messageType=", messageType)
	fmt.Fprintln(writer, "pendingNonce=",state)
	fmt.Fprintln(writer, "currentNonce=",info)
	fmt.Fprintln(writer, "eventId=",interceptor.eventId)
	writer.Flush()

	interceptor.commitEvent()
	return interceptor
	//ackFileName := interceptor.fileDir+"/ack/"+interceptor.filename
	//go interceptor.Wait(ackFileName)
	//interceptor.WaitAck(ackFileName)
}

func UpdateIntercept(sendNode int, recvNode int, messageType string, state int, info int) {
	interceptor := &Interceptor{
		sendNode:		sendNode,
		messageType:    messageType,
		state:			state,
		info:			info}

	interceptor.eventId = interceptor.gethash()
	interceptor.fileDir = "/tmp/ipc"
	interceptor.filename = "ethUpdate-"+strconv.Itoa(sendNode)

	//start intercepting messages into the file inside /new/ folder
	newFileName := interceptor.fileDir+"/new/"+interceptor.filename
	newFile, err := os.Create(newFileName) // Truncates if file already exists, be careful!
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	writer := bufio.NewWriter(newFile)
	defer newFile.Close()
	fmt.Fprintln(writer, "sendNode=",sendNode)
	fmt.Fprintln(writer,"messageType=", messageType)
	fmt.Fprintln(writer, "pendingNonce=",state)
	fmt.Fprintln(writer, "currentNonce=",info)
	fmt.Fprintln(writer, "eventId=",interceptor.eventId)
	writer.Flush()

	interceptor.commitEvent()
}


func (self *Interceptor) commitEvent()  {
	glog.V(logger.Info).Infoln("@RD mv from new to send ",self.filename)
	cmd := exec.Command("mv", self.fileDir+"/new/"+self.filename, self.fileDir+"/send/"+self.filename)
	cmd.Run()
}

func (self *Interceptor) Wait(ackFileName string)  {
	glog.V(logger.Info).Infoln("@RD wait() ", ackFileName)
	for  {
		_, err := os.Stat(ackFileName)
		if err==nil {
			self.Exist<-true
			return
		}
		time.Sleep(time.Millisecond*250)
	}
}

func (self *Interceptor) WaitAck(ackFileName string)  {
	glog.V(logger.Info).Infoln("@RD waitAck ", ackFileName)

	//open and read the ack file
	file, err1 := os.Open(ackFileName)
	if err1!=nil {
		panic(err1)
	}
	content := make([]byte, 10)
	file.Read(content)
	file.Close()
	glog.V(logger.Info).Infoln("@RD ackFile : %s ",content)

	//remove the ack file
	cmd := exec.Command("rm", ackFileName)
	cmd.Run()

}

func (self *Interceptor) createFilename() string {
	var filename string
	now := time.Now()
	nanos := now.Nanosecond()
	millis := nanos / 1000000
	filename = "eth-" + strconv.Itoa(self.eventId)+ strconv.Itoa(millis)
	return filename
}


func (self *Interceptor) gethash() int {
	var prime int = 17
	var hash int = 1
	//sendNode, _:= strconv.Atoi(self.sendNode) how to convert string to int
	hash = prime*hash + self.state
	hash = prime*hash + self.info
	hash = prime*hash + self.sendNode
	hash = prime*hash +  self.recvNode

	msgtypeI, _ := strconv.Atoi(self.messageType)
	hash = prime*hash +  msgtypeI

	glog.V(logger.Info).Infoln("@RD hashof > state: ",self.state, "info: ", self.info,
		"sendNode: ", self.sendNode,"recvNode: ", self.recvNode, "messageType: ", msgtypeI)

	return hash
}

func (self *Interceptor) GetAckFileName() string  {
	return self.fileDir+"/ack/"+self.filename
}
