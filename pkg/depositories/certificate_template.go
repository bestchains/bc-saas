/*
Copyright 2023 The Bestchains Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package depositories

const (
	CertificateTemplate = `{
    "image": "./crt_tpl.jpg",
    "locations": [
        {
            "text": "Certificate Info:",
            "x": 130,
            "y": 275,
            "style": "",
            "size": 16
        },
        {
            "text": "Name: %s",
            "inputs": [
                "name"
            ],
            "x": 180,
            "y": 320,
            "style": "",
            "size": 12
        },
        {
            "text": "Owner: %s",
            "inputs": [
                "owner"
            ],
            "x": 180,
            "y": 360,
            "style": "",
            "size": 12
        },
        {
            "text": "ID: %s",
            "inputs": [
                "kid"
            ],
            "x": 180,
            "y": 400,
            "style": "",
            "size": 12
        },
        {
            "text": "Hash: %s",
            "inputs": [
                "contentID"
            ],
            "x": 180,
            "y": 440,
            "style": "",
            "size": 12
        },
        {
            "text": "Transaction Hash: %s",
            "inputs": [
                "transactionHash"
            ],
            "x": 180,
            "y": 480,
            "style": "",
            "size": 12
        },
        {
            "text": "Date: %s",
            "inputs": [
                "trustedTimestamp"
            ],
            "x": 180,
            "y": 520,
            "style": "",
            "size": 12
        },
        {
            "text": "I promise that the content I upload does not violate any laws or regulations or infringe upon the rights or interests of others. \n I understand and agree that I will be solely responsible for any legal consequences and corresponding legal liabilities arising from the content I upload。",
            "x": 90,
            "y": 560,
            "style": "",
            "size": 12
        },
        {
            "text": "Issued At: %s",
            "inputs": [
                "currentDate"
            ],
            "x": 400,
            "y": 680,
            "style": "",
            "size": 12
        }
    ]
}`

	CertificateTemplateCN = `{
    "image": "./crt_tpl.jpg",
    "locations": [
        {
            "text": "存证信息如下:",
            "x": 130,
            "y": 275,
            "style": "",
            "size": 16
        },
        {
            "text": "名称: %s",
            "inputs": [
                "name"
            ],
            "x": 180,
            "y": 320,
            "style": "",
            "size": 12
        },
        {
            "text": "持有人: %s",
            "inputs": [
                "owner"
            ],
            "x": 180,
            "y": 360,
            "style": "",
            "size": 12
        },
        {
            "text": "唯一存证编号: %s",
            "inputs": [
                "kid"
            ],
            "x": 180,
            "y": 400,
            "style": "",
            "size": 12
        },
        {
            "text": "存证哈希: %s",
            "inputs": [
                "contentID"
            ],
            "x": 180,
            "y": 440,
            "style": "",
            "size": 12
        },
        {
            "text": "交易哈希: %s",
            "inputs": [
                "transactionHash"
            ],
            "x": 180,
            "y": 480,
            "style": "",
            "size": 12
        },
        {
            "text": "上链日期: %s",
            "inputs": [
                "trustedTimestamp"
            ],
            "x": 180,
            "y": 520,
            "style": "",
            "size": 12
        },
        {
            "text": "承诺上传内容不存在任何违反法律法规或侵犯他人权利或权益的情况，理解并同意对其上传的内容自行承担因此产生的一切法律后果及相应法律责任。",
            "x": 90,
            "y": 560,
            "style": "",
            "size": 12
        },
        {
            "text": "颁发日期: %s",
            "inputs": [
                "currentDate"
            ],
            "x": 400,
            "y": 680,
            "style": "",
            "size": 12
        }
    ]
}`
)
