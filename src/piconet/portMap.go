package piconet

/*
type EconetPort struct {
	Value       byte
	Description string
}

func NewPort(value byte) (EconetPort, error) {

	var p = EconetPort{
		Value:       value,
		Description: PortMap[value],
	}
	return p, nil
}
*/

var PortMap = map[byte]string{
	0x00: "Immediate Operation",
	0x4D: "MUGINS",
	0x54: "DigitalServicesTapeStore",
	0x90: "FileServerReply",
	0x91: "FileServerData",
	0x93: "Remote",
	0x99: "FileServerCommand",
	0x9C: "Bridge",
	0x9D: "ResourceLocator",
	0x9E: "PrinterServerEnquiryReply",
	0x9F: "PrinterServerEnquiry",
	0xA0: "SJ Research *FAST protocol",
	0xAF: "SJ Research Nexus net find reply port - SJVirtualEconet",
	0xB0: "FindServer",
	0xB1: "FindServerReply",
	0xB2: "TeletextServerCommand",
	0xB3: "TeletextServerPage",
	0xB4: "Teletext",
	0xB5: "Teletext",
	0xD0: "PrinterServerReply",
	0xD1: "PrinterServerData",
	0xD2: "TCPIPProtocolSuite - IP over Econet",
	0xD3: "SIDFrameSlave, FastFS_Control",
	0xD4: "Scrollarama",
	0xD5: "Phone",
	0xD6: "BroadcastControl",
	0xD7: "BroadcastData",
	0xD8: "ImpressionLicenceChecker",
	0xD9: "DigitalServicesSquirrel",
	0xDA: "SIDSecondary, FastFS_Data",
	0xDB: "DigitalServicesSquirrel2",
	0xDC: "DataDistributionControl, Cambridge Systems Design",
	0xDD: "DataDistributionData, Cambridge Systems Design",
	0xDE: "ClassROM, Oak Solutions",
	0xDF: "PrinterSpoolerCommand, Oak Solutions",
	0xE0: "DigitalServicesNetGain1, David Faulkner, Digital Services",
	0xE1: "DigitalServicesNetGain2, David Faulkner, Digital Services",
	0xE2: "AppFS1, Les Want, AppFS",
	0xE3: "AppFS2, Les Want, AppFS",
	0xE4: "AtomWideFaxNet, Martin Coulson / Chris Ross",
	0xE5: "AtomWidePrintNet, Martin Coulson / Chris Ross",
	0xE6: "IotaDataPower, Neil Raine, Iota",
	0xE7: "CDNetServerBroadcast, Ellis Hall, PEP Associates",
	0xE8: "CDNetServerReplies, Ellis Hall, PEP Associates",
	0xE9: "ClassFS_Server, Oak Solutions",
	0xEA: "DigitalServicesTapeStore2, New allocation to replace 0x54",
	0xEB: "DeveloperSupport, Mark/Jon communication port",
	0xEC: "LLS_Net, Longman Logotron S-Net server",
}
