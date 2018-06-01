package store

// Use This file to change Product Info
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegistration struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	PassConfirm string `json:"pass_confirm"`
	Email       string `json:"email"`
	Type        string `json:"type"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

type SupportOption struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Delivery     bool    `json:"delivery"`
	ExpectedDate uint64  `json:"expected_date"`
	Price        float32 `json:"price"`
}

type SupportOptions []SupportOption

type Idea struct {
	Username           string         `json:"username"`
	Title              string         `json:"title"`
	Image              string         `json:"image"`
	Description        string         `json:"description"`
	TotalFundsRequired uint64         `json:"total_funds_required"`
	TotalFundsRaised   uint64         `json:"total_funds_raised"`
	DatePosted         int64          `json:"date_posted"`
	DateEnd            int64          `json:"date_end"`
	NumContributors    uint32         `json:"num_contributors"`
	Beneficiary        string         `json:"beneficiary"`
	Catagory           string         `json:"category"`
	Summery            string         `json:"summary"`
	IdeaSupportOptions SupportOptions `json:"support_options"`
}

type Ideas []Idea
