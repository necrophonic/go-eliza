package eliza

import (
	"testing"

	"github.com/necrophonic/log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAnalysing(t *testing.T) {

	// Override any logging to off during testing
	log.Init(log.LevelNone)

	var r string
	var rb []byte

	Convey("Ask some questions", t, func() {
		r, _ = AnalyseString("Sorry")
		So(r, ShouldEqual, "Please don't apologise.")

		r, _ = AnalyseString("Gobbledigook")
		So(r, ShouldEqual, "I'm not sure I understand you fully.")
	})

	Convey("Performs substitutions", t, func() {
		r, _ = AnalyseString("But I remember your sheep")
		So(r, ShouldEqual, "Do you often think of my sheep ?")
	})

	Convey("Asking the same question should cycle the response", t, func() {
		r, _ = AnalyseString("Sorry")
		So(r, ShouldEqual, "Apologies are not necessary.")

		r, _ = AnalyseString("Sorry")
		So(r, ShouldEqual, "I've told you that apologies are not required.")

		// Loop back around
		r, _ = AnalyseString("Sorry")
		So(r, ShouldEqual, "Please don't apologise.")
	})

	Convey("Analyse a byte array", t, func() {
		rb, _ = Analyse([]byte{'s', 'o', 'r', 'r', 'y'})
		So(string(rb), ShouldEqual, "Apologies are not necessary.")
	})
}
