package main

import (
  "container/list"
  "lrsms_util"
  "time"
  "strconv"
  "go-coap"
  "math/rand"
  "fmt"
  "os"
  "log"
  "math"
)

const(
  fileName string = "exp2_2.csv"
  localhost string = "localhost"
  initialPort int = 5700
  deviceNum int = 1
  //appIDNum int = 1
  iterations int = 30
  minResourceNum int = 51
  ResourceNumIncrement int = 50
  MaxResourceNum int = 501
)

var ConnectedDevices *list.List
var UnConnectedDevices *list.List
var Devices []*lrsms_util.Device
var AppIDs []string
var Resources []*lrsms_util.Resource

func main() {
  time.Sleep(5 * time.Second)
  //start Profile
  //defer profile.Start().Stop()
  //defer profile.Start(profile.MemProfile).Stop()

  //initail out put file
  f, _ := os.Create(fileName)
  defer f.Close()
  f.WriteString(", MinRunTime, AverageRunTime, MaxRuntime\n")

  //init Devices
  initi(deviceNum)

  // run it
  for resourceNum := minResourceNum; resourceNum <= MaxResourceNum; resourceNum += ResourceNumIncrement{
    //init AppIDs
    AppIDs = make([]string, resourceNum)
    for j := 0; j < resourceNum; j++ {
      AppIDs[j] = "app"+strconv.Itoa(j+1)
      Devices[0].AddApp(AppIDs[j])
    }

    Resources = make([]*lrsms_util.Resource, resourceNum*2-1)
    resourceUpdateChannel := make (chan interface{}, resourceNum*2-1)
    Devices[0].Channel = resourceUpdateChannel

    //init resources
    for i := 0; i < resourceNum; i++{
      if i != 0 {
        CacheResourceID := "Resource"+strconv.Itoa(i)
        CacheDependedRes := list.New()
        CacheDependedResID := Resources[0].URI
        CacheDependedRes.PushBack(CacheDependedResID)
        Resources[(2*i-1)] = initResource(Devices[0].AppServerPort, AppIDs[i],
           CacheResourceID, CacheDependedRes)
        Devices[0].CreateResource(AppIDs[i], Resources[(2*i-1)])
      }

      resourceID := "Resource"+strconv.Itoa(i+1)
      dependedRes := list.New()
      if i != 0 {
        dependedResID := localhost+Devices[0].AppServerPort +"/"+AppIDs[i]+"/"+"Resource"+strconv.Itoa(i)
        dependedRes.PushBack(dependedResID)
      }
      Resources[2*i] = initResource(Devices[0].AppServerPort, AppIDs[i], resourceID,
         dependedRes)

      //Add resource to device/apps
      Devices[0].CreateResource(AppIDs[i], Resources[2*i])
    }

    //printLRSMS(Devices[0].LRSMSServerPort)

    //iterations
    var minRunTime time.Duration = math.MaxInt64
    var MaxRunTime time.Duration = 0
    var TotalRunTime time.Duration = 0
    for i2 := 0; i2 < iterations; i2++ {
      startTime := time.Now()
      //update first resource and see how long it take to reach last resource
      Devices[0].UpdateResource(AppIDs[0], Resources[0].URI)
      endTime := time.Now()
      for i3 := 0; i3 < resourceNum*2-1; i3++ {
        MychanelTtem := <- resourceUpdateChannel
    		endTime = MychanelTtem.(time.Time)
    	}
      duration := endTime.Sub(startTime)
      if duration < minRunTime {
        minRunTime = duration
      }
      if duration > MaxRunTime {
        MaxRunTime = duration
      }
      TotalRunTime += duration
      log.Printf("%d itertaion %d done", resourceNum, (i2+1))
      time.Sleep(200 * time.Millisecond)
    }
    AverageRunTime  := int64(TotalRunTime)/int64(iterations)
    minRunTime /= 1000000
    MaxRunTime /= 1000000
    AverageRunTime /= 1000000
    f.WriteString(strconv.Itoa(resourceNum) + ", " +
     strconv.FormatInt(int64(minRunTime), 10) + ", " +
     strconv.FormatInt(int64(AverageRunTime), 10) + ", " +
     strconv.FormatInt(int64(MaxRunTime), 10)+"\n")
  }
  f.Sync()
  //****************************************************************************
  //Halt
  //****************************************************************************
  //hault til user input e
  for input := "";input!="e";{
    fmt.Println("Enter 'e' to terminate")
    fmt.Scanf("%s", &input)
    //fmt.Println("Hello", name)
  }
}

func initi(deviceNum int){
  ConnectedDevices = list.New()
  UnConnectedDevices = list.New()
  Devices = make([]*lrsms_util.Device, deviceNum)

  i := 0
  //init Devices
  for i = 0; i < deviceNum; i++ {
    deviceAppServPort := initialPort + i*2
    deviceLRSNSServPort := initialPort + i * 2 + 1
    deviceAppServPortString := ":" + strconv.Itoa(deviceAppServPort)
    deviceLRSNSServPortString := ":" + strconv.Itoa(deviceLRSNSServPort)
    Devices[i] = lrsms_util.StartDevice(deviceAppServPortString,
      deviceLRSNSServPortString)
  }
}

//init resoruce
func initResource(deviceAppServPort string, appID string,
  resourceID string, dependedRes *list.List)*lrsms_util.Resource{
  resourceURI := localhost+deviceAppServPort+"/"+appID+"/"+resourceID
  createTime := time.Now().Format(time.RFC3339)
  content := []byte(resourceURI + "created at " + createTime)
  resource :=  lrsms_util.NewResource(resourceURI, dependedRes, createTime,
    content)
  return resource
}

func printLRSMS(port string){
  req := coap.Message{
    Type:      coap.Confirmable,
    Code:      coap.Content,
  	MessageID: uint16(rand.Intn(10000)),
  	Payload:   []byte("{}"),
  }
  req.SetPathString("/Dev")
  c, err := coap.Dial("udp", "localhost"+port)
  if err != nil {
  	log.Fatalf("Error dialing: %v", err)
  }
  rv, err := c.Send(req)
  if err != nil {
  	log.Fatalf("Error sending request: %v", err)
  }
  if rv != nil {
  	//log.Printf("Response payload: %s", rv.Payload)
  }
}
