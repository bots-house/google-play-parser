package models

type DataSafety struct {
	SharedData       []DataSafetyData     `json:"shared_data"`
	CollectedData    []DataSafetyData     `json:"collected_data"`
	PrivacyPolicyURL string               `json:"privacy_policy_url"`
	SecurityPractice []DataSafetyPractice `json:"security_practice"`
}

type DataSafetyData struct {
	Data     string `json:"data"`
	Optional bool   `json:"optional"`
	Purpose  string `json:"purpose"`
	Type     string `json:"type"`
}

type DataSafetyPractice struct {
	Description string `json:"description"`
	Practice    string `json:"practice"`
}
