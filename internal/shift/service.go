package shift




//get shift from start data to end date


type ShiftService struct {
	shiftRepo IShiftRepo
}

func NewShiftService(shiftRepo IShiftRepo) *ShiftService {
	return &ShiftService{shiftRepo: shiftRepo}
}

func(ss *ShiftService)RegisterShift(shift *Shift) error {
	return ss.shiftRepo.RegisterShift(shift)
}