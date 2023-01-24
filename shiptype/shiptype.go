package shiptype

type ShipType int

const (
	NotAvailable ShipType = iota
	ReservedForFutureUse1
	ReservedForFutureUse2
	ReservedForFutureUse3
	ReservedForFutureUse4
	ReservedForFutureUse5
	ReservedForFutureUse6
	ReservedForFutureUse7
	ReservedForFutureUse8
	ReservedForFutureUse9
	ReservedForFutureUse10
	ReservedForFutureUse11
	ReservedForFutureUse12
	ReservedForFutureUse13
	ReservedForFutureUse14
	ReservedForFutureUse15
	ReservedForFutureUse16
	ReservedForFutureUse17
	ReservedForFutureUse18
	ReservedForFutureUse19
	WingInGroundAllShips
	WingInGroundHazardousCategoryA
	WingInGroundHazardousCategoryB
	WingInGroundHazardousCategoryC
	WingInGroundHazardousCategoryD
	WingInGroundReservedForFutureUse1
	WingInGroundReservedForFutureUse2
	WingInGroundReservedForFutureUse3
	WingInGroundReservedForFutureUse4
	WingInGroundReservedForFutureUse5
	Fishing
	Towing
	TowingLengthExceeds200MOrBreadthExceeds25M
	DredgingOrUnderwaterOps
	DivingOps
	MilitaryOps
	Sailing
	PleasureCraft
	Reserved1
	Reserved2
	HighSpeedCraftAllShips
	HighSpeedCraftHazardousCategoryA
	HighSpeedCraftHazardousCategoryB
	HighSpeedCraftHazardousCategoryC
	HighSpeedCraftHazardousCategoryD
	HighSpeedCraftReservedForFutureUse1
	HighSpeedCraftReservedForFutureUse2
	HighSpeedCraftReservedForFutureUse3
	HighSpeedCraftReservedForFutureUse4
	HighSpeedCraftNoAdditionalInformation
	PilotVessel
	SearchAndRescueVessel
	Tug
	PortTender
	AntiPollutionEquipment
	LawEnforcement
	SpareLocalVessel1
	SpareLocalVessel2
	MedicalTransport
	NoncombatantShip
	PassengerAllShips
	PassengerHazardousCategoryA
	PassengerHazardousCategoryB
	PassengerHazardousCategoryC
	PassengerHazardousCategoryD
	PassengerReservedForFutureUse1
	PassengerReservedForFutureUse2
	PassengerReservedForFutureUse3
	PassengerReservedForFutureUse4
	PassengerNoAdditionalInformation
	CargoAllShips
	CargoHazardousCategoryA
	CargoHazardousCategoryB
	CargoHazardousCategoryC
	CargoHazardousCategoryD
	CargoReservedForFutureUse1
	CargoReservedForFutureUse2
	CargoReservedForFutureUse3
	CargoReservedForFutureUse4
	CargoNoAdditionalInformation
	TankerAllShips
	TankerHazardousCategoryA
	TankerHazardousCategoryB
	TankerHazardousCategoryC
	TankerHazardousCategoryD
	TankerReservedForFutureUse1
	TankerReservedForFutureUse2
	TankerReservedForFutureUse3
	TankerReservedForFutureUse4
	TankerNoAdditionalInformation
	OtherTypeAllShips
	OtherTypeHazardousCategoryA
	OtherTypeHazardousCategoryB
	OtherTypeHazardousCategoryC
	OtherTypeHazardousCategoryD
	OtherTypeReservedForFutureUse1
	OtherTypeReservedForFutureUse2
	OtherTypeReservedForFutureUse3
	OtherTypeReservedForFutureUse4
	OtherTypeNoAdditionalInformation
)

func (s ShipType) Description() string {
	var desc string

	switch s {
	case ShipType(0):
		desc = "Not available (default)"
	case ShipType(1):
		fallthrough
	case ShipType(2):
		fallthrough
	case ShipType(3):
		fallthrough
	case ShipType(4):
		fallthrough
	case ShipType(5):
		fallthrough
	case ShipType(6):
		fallthrough
	case ShipType(7):
		fallthrough
	case ShipType(8):
		fallthrough
	case ShipType(9):
		fallthrough
	case ShipType(10):
		fallthrough
	case ShipType(11):
		fallthrough
	case ShipType(12):
		fallthrough
	case ShipType(13):
		fallthrough
	case ShipType(14):
		fallthrough
	case ShipType(15):
		fallthrough
	case ShipType(16):
		fallthrough
	case ShipType(17):
		fallthrough
	case ShipType(18):
		fallthrough
	case ShipType(19):
		desc = "Reserved for future use"
	case ShipType(20):
		desc = "Wing in ground (WIG), all ships of this type"
	case ShipType(21):
		desc = "Wing in ground (WIG), Hazardous category A"
	case ShipType(22):
		desc = "Wing in ground (WIG), Hazardous category B"
	case ShipType(23):
		desc = "Wing in ground (WIG), Hazardous category C"
	case ShipType(24):
		desc = "Wing in ground (WIG), Hazardous category D"
	case ShipType(25):
		desc = "Wing in ground (WIG), Reserved for future use"
	case ShipType(26):
		desc = "Wing in ground (WIG), Reserved for future use"
	case ShipType(27):
		desc = "Wing in ground (WIG), Reserved for future use"
	case ShipType(28):
		desc = "Wing in ground (WIG), Reserved for future use"
	case ShipType(29):
		desc = "Wing in ground (WIG), Reserved for future use"
	case ShipType(30):
		desc = "Fishing"
	case ShipType(31):
		desc = "Towing"
	case ShipType(32):
		desc = "Towing: length exceeds 200m or breadth exceeds 25m"
	case ShipType(33):
		desc = "Dredging or underwater ops"
	case ShipType(34):
		desc = "Diving ops"
	case ShipType(35):
		desc = "Military ops"
	case ShipType(36):
		desc = "Sailing"
	case ShipType(37):
		desc = "Pleasure Craft"
	case ShipType(38):
		desc = "Reserved"
	case ShipType(39):
		desc = "Reserved"
	case ShipType(40):
		desc = "High speed craft (HSC), all ships of this type"
	case ShipType(41):
		desc = "High speed craft (HSC), Hazardous category A"
	case ShipType(42):
		desc = "High speed craft (HSC), Hazardous category B"
	case ShipType(43):
		desc = "High speed craft (HSC), Hazardous category C"
	case ShipType(44):
		desc = "High speed craft (HSC), Hazardous category D"
	case ShipType(45):
		desc = "High speed craft (HSC), Reserved for future use"
	case ShipType(46):
		desc = "High speed craft (HSC), Reserved for future use"
	case ShipType(47):
		desc = "High speed craft (HSC), Reserved for future use"
	case ShipType(48):
		desc = "High speed craft (HSC), Reserved for future use"
	case ShipType(49):
		desc = "High speed craft (HSC), No additional information"
	case ShipType(50):
		desc = "Pilot AISOutput"
	case ShipType(51):
		desc = "Search and Rescue vessel"
	case ShipType(52):
		desc = "Tug"
	case ShipType(53):
		desc = "Port Tender"
	case ShipType(54):
		desc = "Anti-pollution equipment"
	case ShipType(55):
		desc = "Law Enforcement"
	case ShipType(56):
		desc = "Spare - Local AISOutput"
	case ShipType(57):
		desc = "Spare - Local AISOutput"
	case ShipType(58):
		desc = "Medical Transport"
	case ShipType(59):
		desc = "Noncombatant ship according to RR Resolution No. 18"
	case ShipType(60):
		desc = "Passenger, all ships of this type"
	case ShipType(61):
		desc = "Passenger, Hazardous category A"
	case ShipType(62):
		desc = "Passenger, Hazardous category B"
	case ShipType(63):
		desc = "Passenger, Hazardous category C"
	case ShipType(64):
		desc = "Passenger, Hazardous category D"
	case ShipType(65):
		desc = "Passenger, Reserved for future use"
	case ShipType(66):
		desc = "Passenger, Reserved for future use"
	case ShipType(67):
		desc = "Passenger, Reserved for future use"
	case ShipType(68):
		desc = "Passenger, Reserved for future use"
	case ShipType(69):
		desc = "Passenger, No additional information"
	case ShipType(70):
		desc = "Cargo, all ships of this type"
	case ShipType(71):
		desc = "Cargo, Hazardous category A"
	case ShipType(72):
		desc = "Cargo, Hazardous category B"
	case ShipType(73):
		desc = "Cargo, Hazardous category C"
	case ShipType(74):
		desc = "Cargo, Hazardous category D"
	case ShipType(75):
		desc = "Cargo, Reserved for future use"
	case ShipType(76):
		desc = "Cargo, Reserved for future use"
	case ShipType(77):
		desc = "Cargo, Reserved for future use"
	case ShipType(78):
		desc = "Cargo, Reserved for future use"
	case ShipType(79):
		desc = "Cargo, No additional information"
	case ShipType(80):
		desc = "Tanker, all ships of this type"
	case ShipType(81):
		desc = "Tanker, Hazardous category A"
	case ShipType(82):
		desc = "jTanker, Hazardous category B"
	case ShipType(83):
		desc = "Tanker, Hazardous category C"
	case ShipType(84):
		desc = "Tanker, Hazardous category D"
	case ShipType(85):
		desc = "Tanker, Reserved for future use"
	case ShipType(86):
		desc = "Tanker, Reserved for future use"
	case ShipType(87):
		desc = "Tanker, Reserved for future use"
	case ShipType(88):
		desc = "Tanker, Reserved for future use"
	case ShipType(89):
		desc = "Tanker, No additional information"
	case ShipType(90):
		desc = "Other Type, all ships of this type"
	case ShipType(91):
		desc = "Other Type, Hazardous category A"
	case ShipType(92):
		desc = "Other Type, Hazardous category B"
	case ShipType(93):
		desc = "Other Type, Hazardous category C"
	case ShipType(94):
		desc = "Other Type, Hazardous category D"
	case ShipType(95):
		desc = "Other Type, Reserved for future use"
	case ShipType(96):
		desc = "Other Type, Reserved for future use"
	case ShipType(97):
		desc = "Other Type, Reserved for future use"
	case ShipType(98):
		desc = "Other Type, Reserved for future use"
	case ShipType(99):
		desc = "Other Type, no additional information"
	default:
		desc = "ERROR: Ship type out of bounds"
	}

	return desc
}
