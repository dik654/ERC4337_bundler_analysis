import { Injectable } from '@nestjs/common';
import * as qrcode from 'qrcode';
import axios from 'axios';

@Injectable()
export class QrcodeService {
    async generateQR(id: string) {
        try {
            const response = await axios.get(`http://localhost:3001/v1/otp/privatekey/${id}`)
            const qrCodeDataURL = await qrcode.toDataURL(response.data.privateKey)
            return `<img src="${qrCodeDataURL}" alt="QR Code" />`
        } catch (error) {
           return error 
        }
    }

    async verifyOtp(id: string, otp: string) {
        try {
            const response = await axios.post(`http://localhost:3001/v1/otp/verify`, {
                id: id,
                otp: otp
            }) 
            return response.data.verification
        } catch (error) {
            return error
        }
    }
}
