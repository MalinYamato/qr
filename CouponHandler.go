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
	"strconv"
)

type CreateCouponsRequest struct {
	Op       string `json:"op"`
	CouponId string `json:"couponId"`
	Name     string `json:"name"`
	Balance  string `json:"balance"`
	Amount   string `json:"amount"`
}
type GetCouponRequest struct {
	Op      string `json:"op"`
	CoupnID string `json:"couponId"`
}
type DeleteCouponRequest struct {
	Op       string `json:"op"`
	CouponID string `json:"couponId"`
}

type PaymentRequest struct {
	Op       string `json:"op"`
	CouponID string `json:"couponId"`
	Amount   string `json:"amount"`
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
		log.Println("Coupon " + _coupons.getAll()[0].FirstName)
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
		sta := r.ParseForm()
		if sta != nil {
			status.Status = ERROR
			status.Detail = "CouponHandler Parseform Err! "
			log.Println("Parse form failed")
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestCreateCoupon)
		if err != nil {
			log.Println("Json decoder error> ", err.Error())
			panic(err)
		}
		coupon.CouponID = requestCreateCoupon.CouponId
		coupon.FirstName = requestCreateCoupon.Name
		coupon.Balance, _ = strconv.Atoi(requestCreateCoupon.Balance)
		coupon.Amount, _ = strconv.Atoi(requestCreateCoupon.Amount)
		log.Println("Creating coupon of " + coupon.FirstName)
		_coupons.Save(coupon)
		status.Status = SUCCESS
	}

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

func GeneralCouponHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	var status Status
	var paymentRequest PaymentRequest
	//var updateCouponBalance UpdateCouponBalance

	log.Println("CouponHandler called")
	status = Status{SUCCESS, ""}
	//defer r.Body.Close()
	var coupon Coupon
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
			log.Println("Json decoder error> ", err.Error())
			panic(err)
		}
		//request.Op = "payment"
		switch request.Op {
		case "payment":
			{
				err := decoder.Decode(&paymentRequest)
				if err != nil {
					log.Println("Json decoder of paymentRequest error> ", err.Error())
					status.Status = ERROR
					//panic(err)
				}
				coupon, _ := _coupons.findCouponByCouponId(paymentRequest.CouponID)
				amount, _ := strconv.Atoi(paymentRequest.Amount)
				coupon.Balance = coupon.Balance - amount
				_coupons.Save(coupon)
				status.Status = SUCCESS
			}
		case "delete":
			{
				var deleteRequest DeleteCouponRequest
				err := decoder.Decode(&deleteRequest)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				_coupons.DeleteById(deleteRequest.CouponID)
				status.Status = SUCCESS
			}
		default:
			{
				log.Println("wrong or non existent Op " + request.Op)
				status.Status = ERROR
				status.Detail = "wrong or non existent Op " + request.Op
			}
		}
	}

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
