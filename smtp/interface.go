package smtp

type ISMTPService interface {
	SendVerificationEmail(to, verificationCode string) error
}
