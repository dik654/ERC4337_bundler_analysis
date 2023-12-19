import { Injectable } from "@nestjs/common";
import { PassportStrategy } from "@nestjs/passport";
import { Strategy } from 'passport-naver-v2';
import { AuthService } from "../auth.service";

@Injectable()
export class NaverStrategy extends PassportStrategy(Strategy) {
    constructor(private authService: AuthService) {
        super({
        clientID: process.env.NAVER_CLIENT_ID,
        clientSecret: process.env.NAVER_CLIENT_SECRET,
        callbackURL: process.env.NAVER_CALLBACK_URL,
        });
    }

    async validate(
        accessToken: string,
        refreshToken: string,
        profile: any,
        done: any,
    ): Promise<any> {
        const user_email = profile._json.email;
        const user_nick = profile._json.nickname;
        const user_provider = profile.provider;
        const user_profile = {
        user_email,
        user_nick,
        user_provider,
        };

        // console.log(accessToken)
        // console.log(refreshToken)
        // console.log(profile)
        // console.log(done)
        // AAAANp2Q--J4SHaN8mLtXpkLhVzjH7G-hJwIpe3ysgK9d79NETx02PHtcnQc5DoYtsrg4C-bDIz9eETdJQuGKAzWrXk
        // EisxW16jOQez0G93rcyMbvsSt1lipWCWisBZ65ao6CHBE8G99gMj78SI0iiZzEUBCHq4dNcu4gisCLR3buqPZsZtgiszWOZYA4kVFVYfa8FAmfvtpZrK1WhyRoyaNTXOCdAsp7
        // {
        //   provider: 'naver',
        //   id: '5IiBsoCNyVl_Ali28i3YxnJzxenj2DwNcASbwkXR6BE',
        //   nickname: undefined,
        //   profileImage: undefined,
        //   age: undefined,
        //   gender: undefined,
        //   email: undefined,
        //   mobile: undefined,
        //   mobileE164: undefined,
        //   name: '홍길동',
        //   birthday: '03-24',
        //   birthYear: '1945',
        //   _raw: '{"resultcode":"00","message":"success","response":{"id":"5IiBsoCNyVl_Ali28i3YxnJzxenj2DwNcASbwkXR6BE","name":"\\uae20\\ub120\\uc115","birthday":"03-24","birthyear":"1945"}}',
        //   _json: {
        //     resultcode: '00',
        //     message: 'success',
        //     response: {
        //       id: '5IiBsoCNyVl_Ali28i3YxnJzxenj2DwNcASbwkXR6BE',
        //       name: '홍길동',
        //       birthday: '03-24',
        //       birthyear: '1945'
        //     }
        //   }
        // }
        // [Function: verified]

        return { accessToken, refreshToken, profile, type: 'once' }
    }
}