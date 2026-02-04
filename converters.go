package goaliniex

import "strings"

func ToAlpha2CountryCode(code string) string {
	const alpha2Length = 2

	if code == "" {
		return ""
	}

	if len(code) == alpha2Length {
		return code
	}

	alpha3ToAlpha2 := map[string]string{
		"AFG": "AF", "ALA": "AX", "ALB": "AL", "DZA": "DZ", "ASM": "AS",
		"AND": "AD", "AGO": "AO", "AIA": "AI", "ATA": "AQ", "ATG": "AG",
		"ARG": "AR", "ARM": "AM", "ABW": "AW", "AUS": "AU", "AUT": "AT",
		"AZE": "AZ", "BHS": "BS", "BHR": "BH", "BGD": "BD", "BRB": "BB",
		"BLR": "BY", "BEL": "BE", "BLZ": "BZ", "BEN": "BJ", "BMU": "BM",
		"BTN": "BT", "BOL": "BO", "BES": "BQ", "BIH": "BA", "BWA": "BW",
		"BVT": "BV", "BRA": "BR", "IOT": "IO", "BRN": "BN", "BGR": "BG",
		"BFA": "BF", "BDI": "BI", "CPV": "CV", "KHM": "KH", "CMR": "CM",
		"CAN": "CA", "CYM": "KY", "CAF": "CF", "TCD": "TD", "CHL": "CL",
		"CHN": "CN", "CXR": "CX", "CCK": "CC", "COL": "CO", "COM": "KM",
		"COD": "CD", "COG": "CG", "COK": "CK", "CRI": "CR", "CIV": "CI",
		"HRV": "HR", "CUB": "CU", "CUW": "CW", "CYP": "CY", "CZE": "CZ",
		"DNK": "DK", "DJI": "DJ", "DMA": "DM", "DOM": "DO", "ECU": "EC",
		"EGY": "EG", "SLV": "SV", "GNQ": "GQ", "ERI": "ER", "EST": "EE",
		"SWZ": "SZ", "ETH": "ET", "FLK": "FK", "FRO": "FO", "FJI": "FJ",
		"FIN": "FI", "FRA": "FR", "GUF": "GF", "PYF": "PF", "ATF": "TF",
		"GAB": "GA", "GMB": "GM", "GEO": "GE", "DEU": "DE", "GHA": "GH",
		"GIB": "GI", "GRC": "GR", "GRL": "GL", "GRD": "GD", "GLP": "GP",
		"GUM": "GU", "GTM": "GT", "GGY": "GG", "GIN": "GN", "GNB": "GW",
		"GUY": "GY", "HTI": "HT", "HMD": "HM", "VAT": "VA", "HND": "HN",
		"HKG": "HK", "HUN": "HU", "ISL": "IS", "IND": "IN", "IDN": "ID",
		"IRN": "IR", "IRQ": "IQ", "IRL": "IE", "IMN": "IM", "ISR": "IL",
		"ITA": "IT", "JAM": "JM", "JPN": "JP", "JEY": "JE", "JOR": "JO",
		"KAZ": "KZ", "KEN": "KE", "KIR": "KI", "PRK": "KP", "KOR": "KR",
		"KWT": "KW", "KGZ": "KG", "LAO": "LA", "LVA": "LV", "LBN": "LB",
		"LSO": "LS", "LBR": "LR", "LBY": "LY", "LIE": "LI", "LTU": "LT",
		"LUX": "LU", "MAC": "MO", "MKD": "MK", "MDG": "MG", "MWI": "MW",
		"MYS": "MY", "MDV": "MV", "MLI": "ML", "MLT": "MT", "MHL": "MH",
		"MTQ": "MQ", "MRT": "MR", "MUS": "MU", "MYT": "YT", "MEX": "MX",
		"FSM": "FM", "MDA": "MD", "MCO": "MC", "MNG": "MN", "MNE": "ME",
		"MSR": "MS", "MAR": "MA", "MOZ": "MZ", "MMR": "MM", "NAM": "NA",
		"NRU": "NR", "NPL": "NP", "NLD": "NL", "NCL": "NC", "NZL": "NZ",
		"NIC": "NI", "NER": "NE", "NGA": "NG", "NIU": "NU", "NFK": "NF",
		"MNP": "MP", "NOR": "NO", "OMN": "OM", "PAK": "PK", "PLW": "PW",
		"PSE": "PS", "PAN": "PA", "PNG": "PG", "PRY": "PY", "PER": "PE",
		"PHL": "PH", "PCN": "PN", "POL": "PL", "PRT": "PT", "PRI": "PR",
		"QAT": "QA", "REU": "RE", "ROU": "RO", "RUS": "RU", "RWA": "RW",
		"BLM": "BL", "SHN": "SH", "KNA": "KN", "LCA": "LC", "MAF": "MF",
		"SPM": "PM", "VCT": "VC", "WSM": "WS", "SMR": "SM", "STP": "ST",
		"SAU": "SA", "SEN": "SN", "SRB": "RS", "SYC": "SC", "SLE": "SL",
		"SGP": "SG", "SXM": "SX", "SVK": "SK", "SVN": "SI", "SLB": "SB",
		"SOM": "SO", "ZAF": "ZA", "SGS": "GS", "SSD": "SS", "ESP": "ES",
		"LKA": "LK", "SDN": "SD", "SUR": "SR", "SJM": "SJ", "SWE": "SE",
		"CHE": "CH", "SYR": "SY", "TWN": "TW", "TJK": "TJ", "TZA": "TZ",
		"THA": "TH", "TLS": "TL", "TGO": "TG", "TKL": "TK", "TON": "TO",
		"TTO": "TT", "TUN": "TN", "TUR": "TR", "TKM": "TM", "TCA": "TC",
		"TUV": "TV", "UGA": "UG", "UKR": "UA", "ARE": "AE", "GBR": "GB",
		"UMI": "UM", "USA": "US", "URY": "UY", "UZB": "UZ", "VUT": "VU",
		"VEN": "VE", "VNM": "VN", "VGB": "VG", "VIR": "VI", "WLF": "WF",
		"ESH": "EH", "YEM": "YE", "ZMB": "ZM", "ZWE": "ZW",
	}

	if alpha2, ok := alpha3ToAlpha2[code]; ok {
		return alpha2
	}

	return code
}

func ParseIDType(documentType string) IDType {
	switch strings.ToUpper(documentType) {
	case "ID_CARD":
		return IDTypeIDCard
	case "PASSPORT":
		return IDTypePassport
	default:
		return IDType(documentType)
	}
}

func ParseGender(genderStr string) Gender {
	switch strings.ToLower(genderStr) {
	case "m", "male":
		return GenderMale
	case "f", "female":
		return GenderFemale
	default:
		return Gender(strings.ToLower(genderStr))
	}
}

func ToPhoneCode(countryCode string) string {
	countryToPhone := map[string]string{
		"AF": "93", "AL": "355", "DZ": "213", "AS": "1", "AD": "376",
		"AO": "244", "AI": "1", "AQ": "672", "AG": "1", "AR": "54",
		"AM": "374", "AW": "297", "AU": "61", "AT": "43", "AZ": "994",
		"BS": "1", "BH": "973", "BD": "880", "BB": "1", "BY": "375",
		"BE": "32", "BZ": "501", "BJ": "229", "BM": "1", "BT": "975",
		"BO": "591", "BA": "387", "BW": "267", "BR": "55", "IO": "246",
		"VG": "1", "BN": "673", "BG": "359", "BF": "226", "BI": "257",
		"KH": "855", "CM": "237", "CA": "1", "CV": "238", "KY": "1",
		"CF": "236", "TD": "235", "CL": "56", "CN": "86", "CX": "61",
		"CC": "61", "CO": "57", "KM": "269", "CK": "682", "CR": "506",
		"HR": "385", "CU": "53", "CW": "599", "CY": "357", "CZ": "420",
		"CD": "243", "DK": "45", "DJ": "253", "DM": "1", "DO": "1",
		"TL": "670", "EC": "593", "EG": "20", "SV": "503", "GQ": "240",
		"ER": "291", "EE": "372", "ET": "251", "FK": "500", "FO": "298",
		"FJ": "679", "FI": "358", "FR": "33", "PF": "689", "GA": "241",
		"GM": "220", "GE": "995", "DE": "49", "GH": "233", "GI": "350",
		"GR": "30", "GL": "299", "GD": "1", "GU": "1", "GT": "502",
		"GG": "44", "GN": "224", "GW": "245", "GY": "592", "HT": "509",
		"HN": "504", "HK": "852", "HU": "36", "IS": "354", "IN": "91",
		"ID": "62", "IR": "98", "IQ": "964", "IE": "353", "IM": "44",
		"IL": "972", "IT": "39", "CI": "225", "JM": "1", "JP": "81",
		"JE": "44", "JO": "962", "KZ": "7", "KE": "254", "KI": "686",
		"XK": "383", "KW": "965", "KG": "996", "LA": "856", "LV": "371",
		"LB": "961", "LS": "266", "LR": "231", "LY": "218", "LI": "423",
		"LT": "370", "LU": "352", "MO": "853", "MK": "389", "MG": "261",
		"MW": "265", "MY": "60", "MV": "960", "ML": "223", "MT": "356",
		"MH": "692", "MR": "222", "MU": "230", "YT": "262", "MX": "52",
		"FM": "691", "MD": "373", "MC": "377", "MN": "976", "ME": "382",
		"MS": "1", "MA": "212", "MZ": "258", "MM": "95", "NA": "264",
		"NR": "674", "NP": "977", "NL": "31", "NC": "687", "NZ": "64",
		"NI": "505", "NE": "227", "NG": "234", "NU": "683", "KP": "850",
		"MP": "1", "NO": "47", "OM": "968", "PK": "92", "PW": "680",
		"PS": "970", "PA": "507", "PG": "675", "PY": "595", "PE": "51",
		"PH": "63", "PN": "64", "PL": "48", "PT": "351", "PR": "1",
		"QA": "974", "CG": "242", "RE": "262", "RO": "40", "RU": "7",
		"RW": "250", "BL": "590", "SH": "290", "KN": "1", "LC": "1",
		"MF": "590", "PM": "508", "VC": "1", "WS": "685", "SM": "378",
		"ST": "239", "SA": "966", "SN": "221", "RS": "381", "SC": "248",
		"SL": "232", "SG": "65", "SX": "1", "SK": "421", "SI": "386",
		"SB": "677", "SO": "252", "ZA": "27", "KR": "82", "SS": "211",
		"ES": "34", "LK": "94", "SD": "249", "SR": "597", "SJ": "47",
		"SZ": "268", "SE": "46", "CH": "41", "SY": "963", "TW": "886",
		"TJ": "992", "TZ": "255", "TH": "66", "TG": "228", "TK": "690",
		"TO": "676", "TT": "1", "TN": "216", "TR": "90", "TM": "993",
		"TC": "1", "TV": "688", "VI": "1", "UG": "256", "UA": "380",
		"AE": "971", "GB": "44", "US": "1", "UY": "598", "UZ": "998",
		"VU": "678", "VA": "379", "VE": "58", "VN": "84", "WF": "681",
		"EH": "212", "YE": "967", "ZM": "260", "ZW": "263",
	}

	if phone, ok := countryToPhone[strings.ToUpper(countryCode)]; ok {
		return phone
	}

	return ""
}

func SplitPhoneNumber(phoneNumber, countryCode string) (string, string) {
	dialCode := ToPhoneCode(countryCode)
	if dialCode == "" {
		return "", phoneNumber
	}

	dialCode = "+" + dialCode

	cleaned := phoneNumber
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	codeWithoutPlus := strings.TrimPrefix(dialCode, "+")

	if local, found := strings.CutPrefix(cleaned, dialCode); found {
		return dialCode, local
	}

	if local, found := strings.CutPrefix(cleaned, codeWithoutPlus); found {
		return dialCode, local
	}

	if local, found := strings.CutPrefix(cleaned, "00"+codeWithoutPlus); found {
		return dialCode, local
	}

	cleaned, _ = strings.CutPrefix(cleaned, "+")

	return dialCode, cleaned
}

func ExtractPhoneDialCode(phoneNumber string) (string, string) {
	cleaned := phoneNumber
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	normalized := cleaned
	if after, found := strings.CutPrefix(normalized, "+"); found {
		normalized = after
	} else if after, found := strings.CutPrefix(normalized, "00"); found {
		normalized = after
	}

	dialCodes := []string{
		// 4-digit codes
		"1684", "1264", "1268", "1242", "1246", "1441", "1284", "1345",
		"1767", "1809", "1829", "1849", "1473", "1671", "1876", "1664",
		"1670", "1787", "1939", "1869", "1758", "1721", "1784", "1868", "1649", "1340",
		// 3-digit codes
		"355", "213", "376", "244", "672", "374", "297", "994",
		"973", "880", "375", "501", "229", "975", "591", "387",
		"267", "246", "673", "359", "226", "257", "855", "237",
		"238", "236", "235", "269", "682", "506", "385", "599",
		"357", "420", "243", "253", "670", "593", "503", "240",
		"291", "372", "251", "500", "298", "679", "358", "689",
		"241", "220", "995", "233", "350", "299", "502", "224",
		"245", "592", "509", "504", "852", "354", "964", "353",
		"972", "225", "962", "254", "686", "383", "965", "996",
		"856", "371", "961", "266", "231", "218", "423", "370",
		"352", "853", "389", "261", "265", "960", "223", "356",
		"692", "222", "230", "262", "691", "373", "377", "976",
		"382", "212", "258", "264", "674", "977", "687", "505",
		"227", "234", "683", "850", "968", "680", "970", "507",
		"675", "595", "351", "974", "242", "250", "590", "290",
		"508", "685", "378", "239", "966", "221", "381", "248",
		"232", "421", "386", "677", "252", "211", "249", "597",
		"268", "963", "886", "992", "255", "228", "690", "676",
		"216", "993", "688", "256", "380", "971", "598", "998",
		"678", "379", "681", "967", "260", "263",
		// 2-digit codes
		"93", "54", "61", "43", "32", "55", "56", "86", "57", "53",
		"45", "20", "33", "49", "30", "36", "91", "62", "98", "39",
		"81", "82", "60", "52", "31", "64", "47", "92", "51", "63",
		"48", "40", "65", "27", "34", "94", "46", "41", "66", "90",
		"44", "58", "84",
		// 1-digit code (NANP - US, Canada, Caribbean)
		"7", "1",
	}

	for _, code := range dialCodes {
		if localNumber, found := strings.CutPrefix(normalized, code); found {
			return "+" + code, localNumber
		}
	}

	return "", cleaned
}
