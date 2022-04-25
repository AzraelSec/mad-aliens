package engine

func ExecStatusString(s ExecutionStatus) string {
	switch s {
	case MAX_ROUND_REACHED:
		return "Max execution round reached"
	case ALL_ALIENS_STUCK:
		return "All alive aliens are stuck"
	case NO_ALIENS_LEFT:
		return "No alive aliens left"
	default:
		return "Unhandled exit status"
	}
}
