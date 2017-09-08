// Copyright 2017 Aleksey Morarash <tuxofil@gmail.com>
//
// Licensed under the BSD 2 Clause License (the "License");
// you may not use the file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://opensource.org/licenses/BSD-2-Clause
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kunaio

func (h History) MinPrice() float64 {
	if len(h) == 0 {
		return 0
	}
	min := h[0].Price
	for i := 1; i < len(h); i++ {
		if h[i].Price < min {
			min = h[i].Price
		}
	}
	return min
}

func (h History) MaxPrice() float64 {
	max := float64(0)
	for _, e := range h {
		if max < e.Price {
			max = e.Price
		}
	}
	return max
}

func (h History) AvgPrice() float64 {
	return h.SumFunds() / h.SumVolume()
}

func (h History) AvgVolume() float64 {
	return h.SumVolume() / float64(len(h))
}

func (h History) AvgFunds() float64 {
	return h.SumFunds() / float64(len(h))
}

func (h History) SumVolume() float64 {
	var sumVolume float64
	for _, e := range h {
		sumVolume += e.Volume
	}
	return sumVolume
}

func (h History) SumFunds() float64 {
	var sumFunds float64
	for _, e := range h {
		sumFunds += e.Funds
	}
	return sumFunds
}

func (h History) MinVolume() float64 {
	if len(h) == 0 {
		return 0
	}
	min := h[0].Volume
	for i := 1; i < len(h); i++ {
		if h[i].Volume < min {
			min = h[i].Volume
		}
	}
	return min
}

func (h History) MaxVolume() float64 {
	max := float64(0)
	for _, e := range h {
		if max < e.Volume {
			max = e.Volume
		}
	}
	return max
}

func (h History) MinFunds() float64 {
	if len(h) == 0 {
		return 0
	}
	min := h[0].Funds
	for i := 1; i < len(h); i++ {
		if h[i].Funds < min {
			min = h[i].Funds
		}
	}
	return min
}

func (h History) MaxFunds() float64 {
	max := float64(0)
	for _, e := range h {
		if max < e.Funds {
			max = e.Funds
		}
	}
	return max
}
