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
  fileName string = "exp4.csv"
  localhost string = "localhost"
  initialPort int = 5700
  baseInitPort int = 9000
  iterations int = 10
  minResourceNum int = 1
  ResourceNumIncrement int = 1
  MaxResourceNum int = 10
)

var ConnectedDevices *list.List
var UnConnectedDevices *list.List
var Devices []*lrsms_util.Device
var AppIDs []string
var Resources []*lrsms_util.Resource

var BaseDevice *lrsms_util.Device
var BaseAppID string
var resourceUpdateChannel chan interface{}
var BaseResources []*lrsms_util.Resource

func main() {
  time.Sleep(5 * time.Second)

  //initail out put file
  f, _ := os.Create(fileName)
  defer f.Close()
  f.WriteString(", MinRunTime, AverageRunTime, MaxRuntime\n")

  //init device connection list
  ConnectedDevices = list.New()
  UnConnectedDevices = list.New()
  Devices = make([]*lrsms_util.Device, MaxResourceNum)
  AppIDs = make([]string, MaxResourceNum)
  Resources = make([]*lrsms_util.Resource, MaxResourceNum)


  // run it
  for resourceNum := minResourceNum; resourceNum <= MaxResourceNum; resourceNum += ResourceNumIncrement{
    //init Devices & init AppIDs
    k := 0
    if resourceNum != minResourceNum{
      k = resourceNum - ResourceNumIncrement
    }
    for ; k < resourceNum; k++ {
      deviceAppServPort := initialPort + k*2
      deviceLRSNSServPort := initialPort + k * 2 + 1
      deviceAppServPortString := ":" + strconv.Itoa(deviceAppServPort)
      deviceLRSNSServPortString := ":" + strconv.Itoa(deviceLRSNSServPort)
      Devices[k] = lrsms_util.StartDevice(deviceAppServPortString,
        deviceLRSNSServPortString)

      AppIDs[k] = "app"+strconv.Itoa(k+1)
      Devices[k].AddApp(AppIDs[k])
    }
    //init resources
    i := 0
    if resourceNum != minResourceNum{
      i = resourceNum - ResourceNumIncrement
    }
    for ; i < resourceNum; i++{
      resourceID := "Resource"+strconv.Itoa(i+1)
      dependedRes := list.New()
      Resources[i] = initResource(Devices[i].AppServerPort, AppIDs[i], resourceID,
         dependedRes)
      //Add resource to device/apps
      Devices[i].CreateResource(AppIDs[i], Resources[i])
      ConnectedDevices.PushBack(Devices[i])
      Devices[i].Connect(ConnectedDevices)
    }

    //init BaseDevice
    baseDeviceAppServPort := baseInitPort + ((resourceNum-minResourceNum)*2)/ResourceNumIncrement
    baseDeviceLRSNSServPort := baseInitPort + 1 + ((resourceNum-minResourceNum)*2)/ResourceNumIncrement
    baseDeviceAppServPortString := ":" + strconv.Itoa(baseDeviceAppServPort)
    baseDeviceLRSNSServPortString := ":" + strconv.Itoa(baseDeviceLRSNSServPort)
    BaseDevice = lrsms_util.StartDevice(baseDeviceAppServPortString,
      baseDeviceLRSNSServPortString)
    BaseAppID = "app0"
    BaseDevice.AddApp(BaseAppID)

    resourceUpdateChannel = make (chan interface{}, resourceNum*2)
    BaseDevice.Channel = resourceUpdateChannel

    time.Sleep(1000 * time.Millisecond)
    //init baseResources
    BaseResources = make([]*lrsms_util.Resource, MaxResourceNum*2)
    for j := 0 ; j < resourceNum; j++{
      CacheResourceID := "Resource"+strconv.Itoa(j+1)
      CacheDependedRes := list.New()
      CacheDependedResID := Resources[j].URI
      CacheDependedRes.PushBack(CacheDependedResID)
      BaseResources[(2*j)] = initResource(BaseDevice.AppServerPort, BaseAppID,
         CacheResourceID, CacheDependedRes)
      BaseDevice.CreateResource(BaseAppID, BaseResources[(2*j)])

      resourceID := "Resource"+strconv.Itoa(j+1+MaxResourceNum)
      dependedRes := list.New()
      dependedResID := localhost+BaseDevice.AppServerPort +"/"+BaseAppID+"/"+"Resource"+strconv.Itoa(j+1)
      dependedRes.PushBack(dependedResID)
      BaseResources[(2*j+1)] = initResource(BaseDevice.AppServerPort, BaseAppID, resourceID,
         dependedRes)
      //Add resource to device/apps
      BaseDevice.CreateResource(BaseAppID, BaseResources[(2*j+1)])
    }


    //iterations
    var minRunTime time.Duration = math.MaxInt64
    var MaxRunTime time.Duration = 0
    var TotalRunTime time.Duration = 0
    for i2 := 0; i2 < iterations; i2++ {
      time.Sleep(1000 * time.Millisecond)

      //update first resource and see how long it take to reach last resource
      for l := 0; l < resourceNum; l++{
        Devices[l].UpdateResource(AppIDs[l], Resources[l].URI)
      }
      startTime := time.Now()
      BaseDeviceE := ConnectedDevices.PushBack(BaseDevice)
      BaseDevice.Connect(ConnectedDevices)
      endTime := time.Now()
      for i3 := 0; i3 < resourceNum*2; i3++ {
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

      time.Sleep(1000 * time.Millisecond)
      ConnectedDevices.Remove(BaseDeviceE)
      BaseDevice.Disconnect(ConnectedDevices)
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
