import { Controller, Get, Req, Res, UseGuards } from '@nestjs/common';
import { Response } from 'express';
import { NaverAuthGuard } from '../auth/guard/naver-auth.guard';

@Controller('user')
export class UserController {
    @UseGuards(NaverAuthGuard)
    @Get('auth/naver')
    async naverlogin() {
        return;
    }

    @UseGuards(NaverAuthGuard)
    @Get('auth/nlogincallback')
    async callback(@Req() req, @Res() res: Response): Promise<any> {
        res.cookie('access_token', req.user.accessToken)
        res.redirect(`http://localhost:8080/qrcode/${req.user.profile.id}`)
        res.end()
    }
}
