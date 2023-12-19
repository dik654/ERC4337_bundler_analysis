"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.QrcodeController = void 0;
const common_1 = require("@nestjs/common");
const qrcode_service_1 = require("./qrcode.service");
let QrcodeController = class QrcodeController {
    constructor(qrcodeService) {
        this.qrcodeService = qrcodeService;
    }
    async getQRcode() {
        return this.qrcodeService.generateQR();
    }
};
__decorate([
    (0, common_1.Get)(),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", Promise)
], QrcodeController.prototype, "getQRcode", null);
QrcodeController = __decorate([
    (0, common_1.Controller)('qrcode'),
    __metadata("design:paramtypes", [qrcode_service_1.QrcodeService])
], QrcodeController);
exports.QrcodeController = QrcodeController;
//# sourceMappingURL=qrcode.controller.js.map