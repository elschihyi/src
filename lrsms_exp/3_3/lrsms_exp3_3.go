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
  fileName string = "exp3_3.csv"
  localhost string = "localhost"
  initialPort int = 5700
  //deviceNum int = 1
  //appIDNum int = 1
  iterations int = 30
  minResourceNum int = 63
  ResourceNumIncrement int = 2
  //MaxResourceNum int = 255
  MaxResourceNum int = 511
)

var ConnectedDevices *list.List
var UnConnectedDevices *list.List
var Devices []*lrsms_util.Device
var resourceUpdateChannels []chan interface{}
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

  //init device connection list
  ConnectedDevices = list.New()
  UnConnectedDevices = list.New()
  Devices = make([]*lrsms_util.Device, MaxResourceNum)
  resourceUpdateChannels = make ([]chan interface{}, MaxResourceNum)
  AppIDs = make([]string, MaxResourceNum)
  Resources = make([]*lrsms_util.Resource, MaxResourceNum*2-1)

  // run it
  for resourceNum := minResourceNum; resourceNum <= MaxResourceNum;
   resourceNum = (resourceNum * ResourceNumIncrement) +1{
    //init Devices & init AppIDs
    k := 0
    if resourceNum != minResourceNum{
      k = (resourceNum - 1) / ResourceNumIncrement
    }

    for ; k < resourceNum; k++ {
      deviceAppServPort := initialPort + k*2
      deviceLRSNSServPort := initialPort + k * 2 + 1
      deviceAppServPortString := ":" + strconv.Itoa(deviceAppServPort)
      deviceLRSNSServPortString := ":" + strconv.Itoa(deviceLRSNSServPort)
      Devices[k] = lrsms_util.StartDevice(deviceAppServPortString,
        deviceLRSNSServPortString)

      if k==0 {
        resourceUpdateChannels[k] = make (chan interface{}, 1)
      } else {
        resourceUpdateChannels[k] = make (chan interface{}, 2)
      }
      Devices[k].Channel = resourceUpdateChannels[k]
      AppIDs[k] = "app"+strconv.Itoa(k+1)
      Devices[k].AddApp(AppIDs[k])
    }

    i := 0
    if resourceNum != minResourceNum{
      i = (resourceNum - 1) / ResourceNumIncrement
    }

    //init resources
    for ; i < resourceNum; i++{
      if i != 0 {
        CacheResourceID := "Resource"+strconv.Itoa(((i+1)/2))
        CacheDependedRes := list.New()
        CacheDependedResID := localhost+Devices[((i-1)/2)].AppServerPort +"/"+AppIDs[((i-1)/2)]+"/"+"Resource"+strconv.Itoa(((i+1)/2))
        CacheDependedRes.PushBack(CacheDependedResID)
        Resources[(2*i-1)] = initResource(Devices[i].AppServerPort, AppIDs[i],
           CacheResourceID, CacheDependedRes)
        Devices[i].CreateResource(AppIDs[i], Resources[(2*i-1)])
      }

      resourceID := "Resource"+strconv.Itoa(i+1)
      dependedRes := list.New()
      if i != 0 {
        dependedResID := localhost+Devices[i].AppServerPort +"/"+AppIDs[i]+"/"+"Resource"+strconv.Itoa(((i+1)/2))
        dependedRes.PushBack(dependedResID)
      }
      Resources[2*i] = initResource(Devices[i].AppServerPort, AppIDs[i], resourceID,
         dependedRes)

      //Add resource to device/apps
      Devices[i].CreateResource(AppIDs[i], Resources[2*i])
      ConnectedDevices.PushBack(Devices[i])
      Devices[i].Connect(ConnectedDevices)
    }

    //iterations
    var minRunTime time.Duration = math.MaxInt64
    var MaxRunTime time.Duration = 0
    var TotalRunTime time.Duration = 0
    for i2 := 0; i2 < iterations; i2++ {
      time.Sleep(2000 * time.Millisecond)
      startTime := time.Now()
      //update first resource and see how long it take to reach last resource
      Devices[0].UpdateResource(AppIDs[0], Resources[0].URI)
      endTime := time.Now()
      for i3 := 0; i3 < resourceNum; i3++ {
        if i3==0 {
          MychanelTtem := <- resourceUpdateChannels[i3]
    		  endTime = MychanelTtem.(time.Time)
        }else{
          MychanelTtem := <- resourceUpdateChannels[i3]
    		  endTime = MychanelTtem.(time.Time)
          MychanelTtem = <- resourceUpdateChannels[i3]
    		  endTime = MychanelTtem.(time.Time)
        }
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
