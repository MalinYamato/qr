//
// Copyright 2017 Malin Yamato --  All rights reserved.
// https://github.com/MalinYamato
//
// MIT License
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Rakuen. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type Payment struct {
	Amount int `json:"amount"`
	Date   int `json:"date"`
	Time   int `json:"time"`
	Long   int `json:"long,omitempty"`
	Lat    int `json:"lat,omitempty"`
}

type Coupon struct {
	CouponID    string             `json:"couponId"`
	Nic         string             `json:"nic,omitempty"`
	Sign        string             `json:"sign,omitempty"`
	Balance     int                `json:"balance"`
	Amount      int                `json:"amount"`
	IssueTime   int64              `json:"issueTime,omitempty"`
	Email       string             `json:"email,omitempty"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName,omitempty"`
	Password    string             `json:"password,omitempty"`
	PictureURL  string             `json:"pictureURL,omitempty"`
	Description string             `json:"description,omitempty"`
	Token       string             `json:"token,omitempty"`
	Payments    map[string]Payment `json:"payments,omitempty"`
	_Coupons    *Coupons
}

/////////////// Person factory ////////////////////

type Coupons struct {
	__coupons map[string]Coupon
}

func (coupons *Coupons) load() {
	if _, err := os.Stat(coupons.path()); err != nil {

		if os.IsNotExist(err) {
			log.Println("The directory: "+coupons.path()+" does not exist, ignore loading", err)
			return
		}
	}
	files, err := ioutil.ReadDir(coupons.path())
	if err != nil {
		log.Println("directory " + coupons.path() + " not found!")
		log.Fatal(err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(coupons.path() + "/" + file.Name() + "/profile.json")
		if err != nil {
			log.Println("file " + coupons.path() + "/" + file.Name() + " not found!")
			log.Fatal(err)
		}
		var coupon Coupon
		err = json.Unmarshal(content, &coupon)
		if err != nil {
			fmt.Println("error:", err)
		}
		coupons.__coupons[coupon.CouponID] = coupon
	}
}

func (c *Coupon) getNic() string {
	if c.Nic == "" {
		return c.FirstName + " " + c.LastName
	} else {
		return c.Nic
	}
}

func (coupons *Coupons) getAll() (Coupons []Coupon) {
	var l = []Coupon{}
	log.Println("getAll")
	for _, c := range coupons.__coupons {
		l = append(l, c)
	}
	log.Println("getAll after")
	if len(coupons.__coupons) > 1 {
		sort.SliceStable(l, func(i, j int) bool {
			return l[i].CouponID < l[j].CouponID
		})
	}
	return l
}

func (coupons *Coupons) findCouponToken(token string) (c Coupon, ok bool) {
	for _, cc := range coupons.__coupons {
		if cc.Token == token {
			return cc, true
		}
	}
	return Coupon{}, false
}

func (coupons *Coupons) findCouponByCouponId(CouponId string) (coupon Coupon, ok bool) {
	coupon, ok = coupons.__coupons[CouponId]
	return
}

func (coupons *Coupons) Add(c Coupon) bool {
	c._Coupons = coupons
	coupons.__coupons[c.CouponID] = c
	return true
}

func (coupons *Coupons) Save(c Coupon) bool {
	coupons.Add(c)

	if _, err := os.Stat(coupons.path()); err != nil {

		if os.IsNotExist(err) {
			log.Println("Creating "+coupons.path(), err)
			path := coupons.path()
			err := os.Mkdir(path, 0777)
			log.Println("Mkdirerr err ", err)
			if err != nil {
				panic(err)
			}
		}
	}

	log.Println("Pass 1 SaveCoupon ")

	path := c.path()
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			log.Println("Mkdirerr err ", err)
			if err != nil {
				panic(err)
			}
		}
	}
	log.Println("Pass 2 SaveCoupon ")

	path = c.path() + "/img"
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0777)
			log.Println("Mkdirerr err ", err)
			if err != nil {
				panic(err)
			}
		}
	}

	log.Println("Pass 3 SaveCoupon ")
	json_coupon, _ := json.Marshal(c)
	err := ioutil.WriteFile(c.path()+"/profile.json", json_coupon, 0777)
	if err != nil {
		panic(err)
	}
	log.Println("Pass 4 SaveCoupon ")
	log.Println("Number of coupons ", len(coupons.__coupons))
	return true
}
func (coupons *Coupons) DeleteById(CouponID string) bool {
	delete(coupons.__coupons, CouponID)
	path := coupons.path()
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if strings.Contains(CouponID, file.Name()) {
			err := os.RemoveAll(path + "/" + file.Name())
			if err != nil {
				log.Println("os error, Could not remove coupoon ")
				return false
			}
		}
	}
	return true
}
func (coupons *Coupons) Delete(coupon Coupon) bool {
	delete(coupons.__coupons, coupon.CouponID)

	return true
}
func (pers *Coupons) path() string {
	return "./coupons"
}

//////////// Person //////////////

func (c *Coupon) path() string {
	return c._Coupons.path() + "/" + string(c.CouponID)
}
