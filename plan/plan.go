package plan

// A Unit is the lowest unit of reading.
type Unit struct {
	Name   string
	Weight int64
	prev   int
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func sessionName(a, b Unit) string {
	if a.Name == b.Name {
		return a.Name
	}
	return a.Name + " through " + b.Name
}

// Plan will generate a reading plan.
func Plan(u []Unit, n int) []Unit {
	if len(u) <= 0 {
		return u
	}

	h, w = len(u), n
	d := make([][]Unit, w)

	for i := 0; i < n; i++ {
		d[i] = make([]Unit, w)
	}

	var running int64
	for j := 0; j < h; j++ {
		running += u[j].Weight
		d[0][j] = Unit{
			Name:   sessionName(u[0], u[j]),
			Weight: running,
			prev:   -1,
		}
	}

	for i := 1; i < w; i++ {
		for j := h - 1; j >= 0; j-- {
			last := d[i-1][j].Weight
			d[i][j].Weight = last
			d[i][j].prev = j

			for k := j - 1; k >= 0; k-- {
				last = d[i-1][k].Weight
				weight := d[0][j].Weight - d[0][k].Weight

				if weight > d[i][j].Weight {
					break
				}
				if max(weight, last) < d[i][j].Weight {
					d[i][j] = Unit{
						Name:   sessionName(u[i], u[j]),
						Weight: max(weight, last),
						prev:   k,
					}
				}
			}
		}
	}

	result := make([]Unit, w)
	j := h - 1
	for i := w - 1; i >= 0 && j >= 0; i-- {
		result[i] = d[i][j]
		j = result[i].prev
	}
	return result
}
