# go Kerberos sasl authtication steps
## 1. go
    go 1.19
## 2. basic software
### 2.1 hbase 1.2.6.1
    https://archive.apache.org/dist/hbase/1.2.6.1/
### 2.2 thrift 0.9.3
    http://archive.apache.org/dist/thrift/0.9.3/
### 2.3 libkrb5
    Ubuntu: sudo apt-get install libkrb5-dev
    MacOS: brew install homebrew/dupes/heimdal --without-x11
    Debian: yum install -y krb5-devel 
## 3. go dependence
### 3.1 go thrift
    go get git.apache.org/thrift.git/lib/go/thrift@0.9.3
### 3.2 go gosasl
    go get -tags kerberos github.com/beltran/gosasl
### 3.3 thrift hbase client
    tar -zxvf hbase-1.2.6.1-src-tar.gz && cd hbase-1.2.6.1
    thrift --out ./ --gen go ./hbase-thrift/src/main/resources/org/apache/hadoop/hbase/thrift/Hbase.thrift
    cp -r ./hbase/* {you application dir}
## 4.export env
    export KRB5CCNAME={your file path}
    export KRB5_CONFIG={your file path}
## 5. install & run & build
    go get github.com/spierman/gosasl
    go run -tags kerberos {XX}.go
    go build -tags kerberos {XX}.go
## 6. example code
    ```
    import (
        "fmt"
        "log"
        "os"
        sasl "github.com/spierman/gosasl"
	"github.com/spierman/gosasl/hbase"
    )

    func main() {
        host := "XX.com"
        port := 9090
        table := "aa"
        connection, err := sasl.Connect(host, port, sasl.WithGSSAPISaslTransport("hbase"))
        if err != nil {
            log.Fatal("Error connecting", err)
            return
        }
        l := log.New(os.Stderr, "GOKRB5 Client: ", log.LstdFlags)
        client := hbase.NewHbaseClientFactory(connection.Transport, connection.ProtocolFactory)
        isExists, err := client.IsTableEnabled([]byte(table))
        if err != nil {
            l.Fatalf("could not load: %v", err)
        }
        fmt.Printf("rst {%s}\n", isExists)
    }
    ```
## 7. issuse
    if the hbase.go report error,please replace like this
    ```
    replace temp  to  string(temp[:])
    
    ```
