package settings

import "testing"

func TestDefaultIsRichByDefault(t *testing.T) {
	d := Default()

	if !d.ShowLog {
		t.Error("ShowLog should default on (rich by default)")
	}
	if !d.MouseCapture {
		t.Error("MouseCapture should default on")
	}
	if d.LogScope != "current" {
		t.Errorf("LogScope = %q, want current", d.LogScope)
	}
	if d.AuditDays != 30 {
		t.Errorf("AuditDays = %d, want 30", d.AuditDays)
	}
	if d.AuditThreshold != 20 {
		t.Errorf("AuditThreshold = %d, want 20", d.AuditThreshold)
	}
}
