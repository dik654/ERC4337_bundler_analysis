"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.QrcodeService = void 0;
const common_1 = require("@nestjs/common");
const qrcode = require("qrcode");
const axios_1 = require("axios");
let QrcodeService = class QrcodeService {
    async generateQR(id) {
        try {
            const response = await axios_1.default.get(`http://localhost:3001/v1/otp/privatekey/${id}`);
            const qrCodeDataURL = await qrcode.toDataURL(response.data.privateKey);
            return `<img src="${qrCodeDataURL}" alt="QR Code" />`;
        }
        catch (error) {
            return error;
        }
    }
    async verifyOtp(id, otp) {
        try {
            const response = await axios_1.default.post(`http://localhost:3001/v1/otp/verify`, {
                id: id,
                otp: otp
            });
            return response.data.verification;
        }
        catch (error) {
            return error;
        }
    }
};
QrcodeService = __decorate([
    (0, common_1.Injectable)()
], QrcodeService);
exports.QrcodeService = QrcodeService;
//# sourceMappingURL=qrcode.service.js.map