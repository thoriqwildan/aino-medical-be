package model

type UserResponseWrapper struct {
	WebResponse[UserResponse]
}

type ErrorWrapper struct {
	WebResponse[any]
}

type TransactionTypeResponseWrapper struct {
	WebResponse[TransactionTypeResponse]
}

type PlanTypeResponseWrapper struct {
	WebResponse[PlanTypeResponse]
}

type LimitationTypeResponseWrapper struct {
	WebResponse[LimitationTypeResponse]
}

type BenefitResponseWrapper struct {
	WebResponse[BenefitResponse]
}

type DepartmentResponseWrapper struct {
	WebResponse[DepartmentResponse]
}