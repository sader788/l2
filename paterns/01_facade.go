package facade

type Client struct {
}

type ClientManager struct {
	clients []Client
}

func (cl ClientManager) FoundClient(id int) Client {
	return Client{}
}

func (cl Client) PushCarNotification(str string) {

}

type Driver struct {
}

type DriverManager struct {
	drivers []Driver
}

func (dm DriverManager) FoundTaxiForClient(cl Client) Driver {
	return Driver{}
}

func (dm Driver) GetCarInfo() string {
	return "car"
}

type Account struct {
}

type AccountManager struct {
	accounts []Account
}

func (am AccountManager) PayToDriver(cl Client, dr Driver) {
}

// фасад
type DriveCaller struct {
	clientMng ClientManager
	accMng    AccountManager
	driverMng DriverManager
}

func NewDriveCaller() DriveCaller {
	return DriveCaller{ClientManager{}, AccountManager{}, DriverManager{}}
}

func (dc DriveCaller) CallTaxi(id int) {
	cl := dc.clientMng.FoundClient(id)
	dr := dc.driverMng.FoundTaxiForClient(cl)

	cl.PushCarNotification(dr.GetCarInfo())

	dc.accMng.PayToDriver(cl, dr)
}

func main() {
	caller := NewDriveCaller()
	caller.CallTaxi(521)
}
