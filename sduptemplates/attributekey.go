package sduptemplates

//AttributeKey is the string identifier of an attribute
type AttributeKey string

const (
	//AttributeActive represents whether the device is currently on or off
	AttributeActive AttributeKey = "active"
	//AttributeColor represents the primary color of the device
	AttributeColor AttributeKey = "color"
	//AttributeTemperature represents the color temperature of the device
	AttributeTemperature AttributeKey = "temperature"
)
