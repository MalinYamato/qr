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

type CreateCouponsRequest struct {
	Op       string `json:"op"`
	CouponId string `json:"couponId"`
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
	Amount   int    `json:"amount"`
}
type GetCouponRequest struct {
	Op       string `json:"op"`
	CouponID string `json:"couponId"`
}
type DeleteCouponRequest struct {
	Op       string `json:"op"`
	CouponID string `json:"couponId"`
}

type PaymentRequest struct {
	Op       string `json:"op"`
	CouponID string `json:"couponId"`
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
}

type Request struct {
	Op string `json:"op"`
}

// Responses

type GetAllCouponsResponse struct {
	Op      string   `json:"op"`
	Status  Status   `json:"status"`
	Coupons []Coupon `json:"coupons"`
}
type GetCouponResponse struct {
	Op     string `json:"op"`
	Status Status `json:"status"`
	Coupon Coupon `json:"coupon"`
}

func GetAllCouponsHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	var allCoupons GetAllCouponsResponse

	log.Println("GetAllCouponHandler called")
	var status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	if r.Method != "POST" {
		status.Status = ERROR
		status.Detail = "CouponHandler wrong HTTP method! " + r.Method
	} else {
		sta := r.ParseForm()
		if sta != nil {
			status.Status = ERROR
			status.Detail = "CouponHandler Parseform Err! "
			log.Println("Parse form failed")
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			log.Println("Fail to Deconde JSON")
			status.Detail = "getAllCoupons decode Err! "
		}
		allCoupons.Coupons = _coupons.getAll()
		allCoupons.Status = status
		json_response, err := json.Marshal(allCoupons)
		if err != nil {
			log.Println("HandlingCoupon json.Marchal returned error %s", err)
			status.Detail = "getAllCoupons parse Err! "
			//	panic(err)
			//	return
		}
		log.Println("GetAllCouponHandler writing back")
		w.Header().Set("Content-Type", "application/json")
		a, err := w.Write(json_response)
		if err != nil {
			log.Println("Handling GetAllCoupon http.write returned error %s %s", err, a)
			//	panic(err)
			//	return
		}
	}
}

func CreateCouponHandler(w http.ResponseWriter, r *http.Request) {
	var requestCreateCoupon CreateCouponsRequest
	var status Status
	//var updateCouponBalance UpdateCouponBalance
	log.Println("CouponHandler called")
	status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	var coupon Coupon
	if r.Method != "POST" {
		status.Status = ERROR
		status.Detail = "CouponHandler wrong HTTP method! " + r.Method
	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestCreateCoupon)
		if err != nil {
			log.Println("Json decoder error> ", err.Error())
			panic(err)
		}
		coupon.CouponID = requestCreateCoupon.CouponId
		coupon.FirstName = requestCreateCoupon.Name
		coupon.Balance = requestCreateCoupon.Balance
		coupon.Amount = requestCreateCoupon.Amount
		log.Println("Creating coupon of " + coupon.FirstName)
		_coupons.Save(coupon)

		json_response, err := json.Marshal(status)
		if err != nil {
			log.Println("HandlingCoupon json.Marchal returned error %s", err)
			panic(err)
			return
		}
		log.Println("CouponHandler writing back status of " + coupon.FirstName)
		w.Header().Set("Content-Type", "application/json")
		a, err := w.Write(json_response)
		if err != nil {
			log.Println("HandlingCoupon http.write returned error %s %s", err, a)
			panic(err)
			return
		}
	}
}

func UpdateCouponHandler(w http.ResponseWriter, r *http.Request) {
	var status Status
	var paymentRequest PaymentRequest
	log.Println("UpdateCouponHandler called")
	status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	if r.Method != "POST" {
		status.Status = ERROR
		status.Detail = "CouponHandler wrong HTTP method! " + r.Method
	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&paymentRequest)
		if err != nil {
			log.Println("Json decoder of paymentRequest error> ", err.Error())
			status.Status = ERROR
			//panic(err)
		}
		coupon, stat := _coupons.findCouponByCouponId(paymentRequest.CouponID)
		if stat == false {
			status.Status = WARNING
			status.Detail = "There are no coupons!"
		} else {
			switch paymentRequest.Op {
			case "addBalance":
				{
					log.Println("balance %d  &d", coupon.Balance, paymentRequest.Balance)
					coupon.Balance = coupon.Balance + paymentRequest.Balance
				}
			case "payment":
				{
					coupon.Balance = coupon.Balance - paymentRequest.Amount
				}
			default:
				{
					log.Println("UpdateCouponHander no such op %s", paymentRequest.Op)
				}
			}
			_coupons.Save(coupon)
			status.Status = SUCCESS
		}
		json_response, err := json.Marshal(status)
		if err != nil {
			log.Println("HandlingCoupon json.Marchal returned error %s", err)
			panic(err)
			return
		}
		log.Println("UpdateCouponHandler writing back status of " + coupon.FirstName)
		w.Header().Set("Content-Type", "application/json")
		a, err := w.Write(json_response)
		if err != nil {
			log.Println("HandlingCoupon http.write returned error %s %s", err, a)
			panic(err)
			return
		}
	}
}

func GetCouponHandler(w http.ResponseWriter, r *http.Request) {
	var status Status
	var getCouponRequest GetCouponRequest
	var getCouponResponse GetCouponResponse
	log.Println("CouponHandler called")
	status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	if r.Method != "POST" {
		status.Status = ERROR
		status.Detail = "CouponHandler wrong HTTP method! " + r.Method
	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&getCouponRequest)
		if err != nil {
			log.Println("Json decoder of paymentRequest error> ", err.Error())
			status.Status = ERROR
			status.Detail = "Json decoder of paymentRequest error> "
			//panic(err)
		}
		coupon, stat := _coupons.findCouponByCouponId(getCouponRequest.CouponID)
		if stat == false {
			status.Status = WARNING
			status.Detail = "There are no coupons!"
		} else {
			getCouponResponse.Coupon = coupon
			status.Status = SUCCESS
		}
		getCouponResponse.Status = status
		json_response, err := json.Marshal(getCouponResponse)
		if err != nil {
			log.Println("HandlingCoupon json.Marchal returned error %s", err)
			panic(err)
			return
		}
		log.Println("CouponHandler writing back status of " + coupon.FirstName)
		w.Header().Set("Content-Type", "application/json")
		a, err := w.Write(json_response)
		if err != nil {
			log.Println("HandlingCoupon http.write returned error %s %s", err, a)
			panic(err)
			return
		}
	}
}
