import { Injectable } from '@nestjs/common';
import { Cron, CronExpression } from '@nestjs/schedule';
import * as qrcode from 'qrcode';
import axios from 'axios';


@Injectable()
export class QrcodeService {
    @Cron(CronExpression.EVERY_5_SECONDS)
    async generateQR() {
        try {
            // const response = await axios.get('http://localhost:8080/qrcode')
            // console.log('Response:', response.data)
            const qrCodeDataURL = await qrcode.toDataURL(Math.random().toString())
            return `<img src="${qrCodeDataURL}" alt="QR Code" />`
        } catch (error) {
           return error 
        }
    }
}
