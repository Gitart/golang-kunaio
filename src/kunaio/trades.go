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

type Trades []Trade

func (t Trades) SumVolume() float64 {
	var s float64
	for _, e := range t {
		s += e.Volume
	}
	return s
}

func (t Trades) SumFunds() float64 {
	var s float64
	for _, e := range t {
		s += e.Funds
	}
	return s
}

func (t Trades) AvgPrice() float64 {
	return t.SumFunds() / t.SumVolume()
}

func (t Trades) AvgVolume() float64 {
	return t.SumVolume() / float64(len(t))
}

func (t Trades) AvgFunds() float64 {
	return t.SumFunds() / float64(len(t))
}
