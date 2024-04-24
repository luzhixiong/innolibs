package utils

var (
	newYorkDeliverableZipCodes = []string{"11004", "11101", "11102", "11103", "11104", "11105", "11106", "11354", "11355", "11356", "11357", "11358", "11360", "11361", "11362", "11363", "11364", "11365", "11366", "11367", "11368", "11369", "11370", "11372", "11373", "11374", "11375", "11377", "11507", "11509", "11510", "11710", "11714", "11514", "11516", "11554", "11732", "11518", "11003", "11735", "11001", "11002", "11003", "11010", "11520", "11530", "11531", "11535", "11542", "11545", "11545", "11547", "11020", "11021", "11022", "11023", "11024", "11025", "11026", "11027", "11548", "11549", "11550", "11040", "11557", "11801", "11802", "11815", "11819", "11854", "11855", "11096", "11558", "11753", "11559", "11756", "11560", "11561", "11563", "11565", "11030", "11758", "11762", "11566", "11765", "11501", "11502", "11514", "11040", "11041", "11042", "11043", "11044", "11099", "11572", "11804", "11568", "11771", "11803", "11050", "11051", "11052", "11053", "11054", "11055", "11570", "11571", "11575", "11576", "11577", "11579", "11783", "11773", "11791", "11553", "11580", "11581", "11582", "11793", "11552", "11568", "11590", "11596", "11797", "11598", "11705", "11709", "11713", "11715", "11716", "11717", "11932", "11718", "11719", "11933", "11934", "11720", "11721", "11722", "11724", "11725", "11726", "11727", "11935", "11729", "11937", "11954", "11730", "11939", "11940", "11731", "11942", "11733", "11941", "11738", "06390", "11739", "11740", "11944", "11946", "11749", "11760", "11741", "11742", "11743", "11746", "11750", "11746", "11747", "11749", "11751", "11752", "11947", "11754", "11755", "11948", "11757", "11949", "11950", "11951", "11952", "11763", "11747", "11775", "11953", "11764", "11954", "11955", "11766", "11767", "11956", "11703", "11768", "11769", "11770", "11957", "11772", "11958", "11777", "11776", "11959", "11960", "11961", "11901", "11970", "11778", "11779", "11963", "11964", "11962", "11780", "11782", "11784", "11964", "11965", "11967", "11786", "11787"}
	newYorkSpecialZipCodes     = []string{"07101", "07102", "07103", "07104", "07105", "07106", "07107", "07108", "07112", "07114", "07302", "07304", "07305", "07306", "07307", "07310", "07311", "07395", "07501", "07502", "07503", "07504", "07505", "07508", "07510", "07201", "07202", "07206", "08817", "08818", "08820", "08837", "07095", "08701", "08753", "08755", "08609", "08610", "08611", "08619", "08608", "08609", "08610", "08611", "08618", "08619", "07011", "07012", "07013", "07014", "08102", "08103", "08104", "08105", "08002", "08003", "08034", "07055", "07087", "07002", "07017", "07018", "08901", "08902", "08861", "07030", "07093", "07060", "07062", "07601", "08872", "07036", "08401", "08402", "08406", "11201", "11203", "11204", "11205", "11206", "11207", "11208", "11209", "11210", "11211", "11212", "11213", "11214", "11215", "11216", "11217", "11218", "11219", "11220", "11221", "11222", "11223", "11224", "11225", "11226", "11228", "11229", "11230", "11231", "11232", "11233", "11234", "11235", "11236", "11237", "11238", "10451", "10452", "10453", "10454", "10455", "10456", "10457", "10458", "10459", "10460", "10461", "10462", "10463", "10464", "10465", "10466", "10467", "10468", "10469", "10470", "10471", "10472", "10473", "10474", "10475", "10001", "10002", "10003", "10004", "10005", "10006", "10007", "10009", "10010", "10011", "10012", "10013", "10014", "10016", "10017", "10018", "10019", "10021", "10022", "10023", "10024", "10025", "10026", "10027", "10028", "10029", "10030", "10031", "10032", "10033", "10034", "10035", "10036", "10037", "10038", "10039", "10040", "10044", "10065", "10069", "10075", "10103", "10110", "10111", "10112", "10115", "10128", "10152", "10153", "10154", "10162", "10301", "10302", "10303", "10304", "10305", "10306", "10307", "10308", "10309", "10310", "10311", "10312", "10314"}
)

func IsNewYorkDeliverableZipCode(zipcode string) bool {
	for _, z := range newYorkDeliverableZipCodes {
		if zipcode == z {
			return true
		}
	}
	return false
}

func IsNewYorkSpecialZipCode(zipcode string) bool {
	for _, z := range newYorkSpecialZipCodes {
		if zipcode == z {
			return true
		}
	}
	return false
}