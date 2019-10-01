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
	"log"
	"net/http"
)

type CouponRequest struct {
	Op      string `json:"op"`
	Value   string `json:"value"`
	CoupnID string `json:"couponID"`
	Name    string `json:"name"`
}
type CouponResponse struct {
	Op       string   `json:"op"`
	Status   Status   `json:"status"`
	Value    string   `json:"value"`
	Balance  int      `json:"value"`
	CouponID string   `json:"couponID"`
	Name     string   `json:"name"`
	Coupons  []Coupon `json:"coupons"`
}

func CouponHandler(w http.ResponseWriter, r *http.Request) {
	var request CouponRequest
	var response CouponResponse
	response.Op = response.Op
	var status Status
	status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	if r.Method != "POST" {
		status.Status = ERROR
		status.Detail = "CouponHandler wrong HTTP method!"
	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			log.Println("ERR> ", err)
		}
		switch request.Op {
		case "deleteCoupon":
			{
			}
		case "addCoupon":
			{
			}
		case "updateBalance":
			{
			}
		case "getCoupon":
			{
			}
		case "getAllCoupons":
			{
			}
		default:
			{
			}
		}
	}
	response.Status = status
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("HandlingCoupon json.Marchal returned error %s", err)
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	a, err := w.Write(data)
	if err != nil {
		log.Println("HandlingCoupon http.write returned error %s", err)
		panic(err)
		return
	}
}