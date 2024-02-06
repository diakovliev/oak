package ray

import (
	"testing"

	"github.com/diakovliev/oak/v4/alg/floatgeom"
	"github.com/diakovliev/oak/v4/alg/span"
	"github.com/diakovliev/oak/v4/collision"
)

func TestEmptyRaycasts(t *testing.T) {
	t.Skip()
	collision.DefaultTree.Clear()
	vRange := span.NewLinear(3.0, 359.0)
	tests := 100
	for i := 0; i < tests; i++ {
		p1 := floatgeom.Point2{vRange.Poll(), vRange.Poll()}
		p2 := floatgeom.Point2{vRange.Poll(), vRange.Poll()}
		if len(Cast(p1, p2)) != 0 {
			t.Fatalf("cast found a point in the empty tree")
		}
		if len(CastTo(p1, p2)) != 0 {
			t.Fatalf("cast to found a point in the empty tree")
		}
		if len(ConeCast(p1, p2)) != 0 {
			t.Fatalf("cone cast found a point in the empty tree")
		}
		if len(ConeCastTo(p1, p2)) != 0 {
			t.Fatalf("cone cast to found a point in the empty tree")
		}
	}
}
