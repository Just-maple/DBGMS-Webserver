# A WebServer Based on Golang Gin and Mgo as Database

### Client :  https://git.asoapp.com/wulingfeng/gin-vue-webclient

### DEMO build : go build -o server ./demo

#### Run order:

#### 0. Init your config  please set a mongodb url as db dial

#### 1. Init Server

#### 2. Init Database      --handler.NewDataBase()

#### 3. Register API       --handler.RegisterAPI()

#### 4. Init Meta Config   --handler.InitMetaConfig()

#### 5. Start Server       --server.Start()

#### you can visit http://localhost:8388/ to start (default port,you can change it in config)