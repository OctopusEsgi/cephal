package gameserver

// DATA SUR LES JEU + VERIF des variable d'env

var internalGamePortsMap = map[string]map[string][]string{
	"mindustryesgi": {
		"tcp": {"6567"},
		"udp": {"6567"},
	},
	"programmeperso": {
		"tcp": {"4000", "4001", "4002", "5892"},
		"udp": {"9899"},
	},
}

var internalGameSpec = map[string]map[string]int{
	"mindustryesgi": {
		"core": 1,
		"ram":  512,
	},
}

// var internalGameEnv := {
// "mindustryesgi":{
// "VERSION": ???,
// "MAP": ???,
// "MODE": ???,
// "PLAYERSLIMIT": ???,
// }
// }
