import { Response } from 'express';
export declare class UserController {
    naverlogin(): Promise<void>;
    callback(req: any, res: Response): Promise<any>;
}
