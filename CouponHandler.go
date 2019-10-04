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
	Op       string `json:"op"`
	Status   Status `json:"status"`
	CouponID string `json:"couponID"`
	Name     string `json:"name"`
}

func CouponHandler(w http.ResponseWriter, r *http.Request) {
	var request CouponRequest
	var response CouponResponse
	var status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	if r.Method == "POST" {
		status.Status = ERROR
		status.Detail = "CouponHandler wrong HTTP method! " + r.Method
	} else {
		sta := r.ParseForm()
		if sta != nil {
			status.Status = ERROR
			status.Detail = "CouponHandler Parseform Err! "
		}

		//response.Op = r.Form.Get("Op")
		status.Detail = r.Method + " " + request.Op
		response.CouponID = r.Method
		response.Name = r.Form.Get("Name")
		request.Op = "addCoupon"
		response.Op = request.Op
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
				log.Println("HandlingCoupon wrong Op %s", response.Op)
				status.Status = ERROR
				status.Detail = "HandlingCoupon wrong Op"
			}
		}
	}

	response.Status = status
	ss, err := json.Marshal(status)
	if err != nil {
		log.Println("HandlingCoupon json.Marchal returned error %s", err)
		panic(err)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	a, err := w.Write(ss)
	if err != nil {
		log.Println("HandlingCoupon http.write returned error %s %s", err, a)
		panic(err)
		return
	}
}
