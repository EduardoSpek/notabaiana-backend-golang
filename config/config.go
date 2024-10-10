package config

var (
	// Suamusica_api_version is the version of the Suamusica API
	Suamusica_UrlSite            = "https://suamusica.com.br"
	Suamusica_api_version string = "1025"

	//Dimensions of the banners
	Banner_pc_dimensions     = [2]int{1300, 190}
	Banner_tablet_dimensions = [2]int{726, 106}
	Banner_mobile_dimensions = [2]int{386, 386}

	//News limit per page
	News_AllowedDomains = "www.bahianoticias.com.br"
	News_LimitPerPage   = 1000
	News_PerPage        = 24

	//Downloads limit
	Downloads_PerPage = 24

	//Dimensions of the Contact
	Contact_ImageDimensions = [2]int{1280, 1280}
)
