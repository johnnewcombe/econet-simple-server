package econet

func ProcessFunctionCode(functionCode byte, command string, srcStationId byte, srcNetworkId byte) []byte {

	reply := CLIDecode(tidyText(command), srcStationId, srcNetworkId)

	return reply
}
