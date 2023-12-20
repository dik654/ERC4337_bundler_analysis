export declare class QrcodeService {
    generateQR(id: string): Promise<any>;
    verifyOtp(id: string, otp: string): Promise<any>;
}
