export interface StakeProvider {
    serviceFee: number;
    apr: number;
    identity: string;
}
export interface Prices {
    egld: string;
    mex: string;
}
export interface StrategyResult {
    ProfitInEgld: string;
    ProfitInMex: string;
    ProfitInUSD: string;
    ROI: string;
    TotalBalanceInEgld: string;
    TotalBalanceInMex: string;
    TotalBalanceInUsd: string;
}
export interface CalculateResponse {
    prices: Prices;
    results: {
        [key: string]: StrategyResult;
    };
}
export interface AppLogicState {
    stakeProviders: StakeProvider[];
    prices: Prices;
    enabledMexAndEgld: boolean;
    enabledMexRewardsLocked: boolean;
}
export declare class AppLogic extends HTMLElement {
    _state: AppLogicState;
    constructor();
    updatePrices(): void;
    displayResults(results: {
        [key: string]: StrategyResult;
    }): void;
    displayResult(name: string, result?: StrategyResult): void;
    fetchStakingProviders(): void;
    fetchPrices(): void;
    fetchSubmit(request: {
        [key: string]: any;
    }): void;
    handleInputEvent(name: string, value: string): void;
    handleSubmit(): void;
    connectedCallback(): void;
}
export declare class RewardDisplay extends HTMLElement {
    result: StrategyResult;
    strategyName: string;
    constructor(strategyName?: string, result?: StrategyResult);
    connectedCallback(): void;
}
