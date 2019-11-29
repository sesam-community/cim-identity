package main

// Backend sets up the http.Client with authentication for backend services
func (s *Server) Backend() error {
	// cfg := s.options.backend
	// if len(cfg.serviceURL) == 0 || len(cfg.clientID) == 0 || len(cfg.clientSecret) == 0 || len(cfg.tokenURL) == 0 {
	// 	s.Logf(logWARN, "Backend services disabled due to missing configuration options. DataIdentity-API services not being used for UUIDs.\n")
	// 	s.client = http.DefaultClient
	// 	return nil
	// }

	// ctx := context.Background()

	// conf := &oauth2.Config{
	// 	ClientID:     "YOUR_CLIENT_ID",
	// 	ClientSecret: "YOUR_CLIENT_SECRET",
	// 	Scopes:       []string{"SCOPE1", "SCOPE2"},
	// 	Endpoint: oauth2.Endpoint{
	// 		TokenURL: "https://provider.com/o/oauth2/token",
	// 		AuthURL:  "https://provider.com/o/oauth2/auth",
	// 	},
	// }
	return nil
}
