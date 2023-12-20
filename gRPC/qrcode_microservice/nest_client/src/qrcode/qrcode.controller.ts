import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { QrcodeService } from './qrcode.service';

@Controller('qrcode')
export class QrcodeController {
    constructor(private readonly qrcodeService: QrcodeService) {}
    @Get(":id")
    async getQRcode(@Param("id") id: string) {
        return this.qrcodeService.generateQR(id);
    }

    @Post('/verify')
    async verifyOtp(@Body() body) {
        return this.qrcodeService.verifyOtp(body.id, body.otp);
    }
}

