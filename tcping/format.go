package tcping

import (
	"fmt"
	"time"
)

func FormatResult(result ProbeResult, useColor bool) string {
	tx := result.TxPacket.Header
	rx := result.RxPacket.Header
	var str string
	str = str + fmt.Sprintf("[      %-15s        %-4dms          %-15s ]\n", result.TxPacket.IP, result.Latency()/int64(time.Millisecond), result.RxPacket.IP)
	str = str + fmt.Sprintf("[ SRC: %5d ] [ DST: %5d ]     [ SRC: %5d ] [ DST: %5d ]\n", tx.Src, tx.Dst, rx.Src, rx.Dst)
	str = str + fmt.Sprintf("[ SEQ: %20d ]     [ SEQ: %20d ]\n", tx.Seq, rx.Seq)
	str = str + fmt.Sprintf("[ ACK: %20d ]     [ ACK: %20d ]\n", tx.Ack, rx.Ack)
	str = str + fmt.Sprintf("[ FLG: %s%s%s%s%s%s] [ WIN: %d ]     [ FLG: %s%s%s%s%s%s] [ WIN: %d ]\n",
		FlagEntry(tx, URG, useColor),
		FlagEntry(tx, ACK, useColor),
		FlagEntry(tx, PSH, useColor),
		FlagEntry(tx, RST, useColor),
		FlagEntry(tx, SYN, useColor),
		FlagEntry(tx, FIN, useColor),
		tx.Window,
		FlagEntry(rx, URG, useColor),
		FlagEntry(rx, ACK, useColor),
		FlagEntry(rx, PSH, useColor),
		FlagEntry(rx, RST, useColor),
		FlagEntry(rx, SYN, useColor),
		FlagEntry(rx, FIN, useColor),
		rx.Window,
	)
	str = str + fmt.Sprintf("[ SUM: %5d ] [ URG: %5d ]     [ SUM: %5d ] [ URG: %5d ]\n", tx.Checksum, tx.Urgent, rx.Checksum, rx.Urgent)
	//for _, o := range tx.Options {
	//	str = str + fmt.Sprintf("[ Option: kind=%d len=%d data=%v ]\n", o.Kind, o.Length, o.Data)
	//}
	return str
}

func FlagEntry(header TCPHeader, flag byte, color bool) string {
	if header.HasFlag(flag) {
		switch flag {
		case URG:
			return "U"
		case ACK:
			return "A"
		case PSH:
			return "P"
		case RST:
			return "R"
		case SYN:
			return "S"
		case FIN:
			return "F"
		default:
			return " "
		}
	} else {
		return "_"
	}
}
