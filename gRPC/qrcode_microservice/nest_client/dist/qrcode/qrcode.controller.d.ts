import { QrcodeService } from './qrcode.service';
export declare class QrcodeController {
    private readonly qrcodeService;
    constructor(qrcodeService: QrcodeService);
    getQRcode(): Promise<any>;
}
