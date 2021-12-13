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
    }
}

export interface AppLogicState {
    stakeProviders: StakeProvider[];
    prices: Prices;

    enabledMexAndEgld: boolean;
    enabledMexRewardsLocked: boolean;
}


export class AppLogic extends HTMLElement {
    _state: AppLogicState;

    constructor() {
        super();

        this._state = {
            stakeProviders: [],
            prices: {
                egld: "0",
                mex: "0",
            },
            enabledMexAndEgld: true,
            enabledMexRewardsLocked: false
        };
    }

    updatePrices() {
        const priceEgldEl = document.querySelector('#section-egld h3 span')! as HTMLSpanElement;
        const priceMexEl = document.querySelector('#section-mex h3 span')! as HTMLSpanElement;

        priceEgldEl.innerText = "$" + this._state.prices.egld;
        priceMexEl.innerText = "$" + this._state.prices.mex;
    }

    displayResults(results: {[key: string]: StrategyResult;}) {
        console.debug(results);

        const rewardsEl = document.querySelectorAll('#results reward-display')!;
        rewardsEl.forEach(el => {
            el.remove();
        });

        this.displayResult('EGLD (Hold in Wallet)', results['egld_hold']);
        this.displayResult('EGLD (Staking)', results['egld_stake']);
        this.displayResult('EGLD (Stake & Redelegate)', results['egld_redelegate']);
        this.displayResult('MEX (Staking)', results['mex_stake']);
        this.displayResult('MEX (Stake & Redelegate)', results['mex_redelegate']);
    }

    displayResult(name: string, result?: StrategyResult) {
        if (!result) {
            return;
        }

        const resultsEl = document.querySelector('#results')! as HTMLDivElement;
        const rewardEl = new RewardDisplay(name, result);

        resultsEl.appendChild(rewardEl);

        console.debug(name, result);
    }


    fetchStakingProviders() {
        const selectProviderEl = document.querySelector('select[name="egld-staking-provider"]')! as HTMLSelectElement;

        fetch("https://istari-api.troscot.com/api/egld_staking_providers")
            .then(res => {
                if (res.status !== 200) {
                    throw new Error("bad status code: " + res.status);
                }

                return res.json();
            })
            .then(json => {
                selectProviderEl.innerHTML = '';

                this._state.stakeProviders = json['staking_providers'];
                this._state.stakeProviders.forEach(provider => {
                    const option = document.createElement('option');
                    option.value = provider.identity;
                    option.innerText = `apr[${provider.apr}] - ${provider.identity}`;
                    selectProviderEl.appendChild(option);
                });

                selectProviderEl.dispatchEvent(new Event('change'));
            })
            .catch(err => {
                alert("failed fetching staking providers: " + err);
                console.error(err);
            });

        return;
    }

    fetchPrices() {
        fetch("https://istari-api.troscot.com/api/prices")
            .then(res => {
                if (res.status != 200) {
                    throw new Error("bad status code: " + res.status);
                }

                return res.json();
            })
            .then(json => {
                this._state.prices = json['prices'];
                this.updatePrices();
            })
            .catch(err => {
                alert("failed fetching staking providers: " + err);
                console.error(err);
            });
    }

    fetchSubmit(request: { [key:string]: any }) {
        const json = JSON.stringify(request);

        fetch("https://istari-api.troscot.com/api/calculate_profit", { method: "POST", body: json })
            .then(res => {
                if (res.status != 200) {
                    throw new Error("bad status code: " + res.status);
                }

                return res.json() as Promise<CalculateResponse>;
            })
            .then(json => {
                console.debug('response', json);

                this._state.prices = json.prices;
                this.updatePrices();
                this.displayResults(json.results);
            })
            .catch(err => {
                alert("failed calculating profit: " + err);
                console.error(err);
            });
    }

    handleInputEvent(name: string, value: string) {
        console.debug('handle input event', {
            'name': name,
            'value': value,
        });
        switch (name) {
            case 'toggle-egld-only':
                this._state.enabledMexAndEgld = !this._state.enabledMexAndEgld;
                const sectionEl = document.getElementById('section-mex')!;
                const buttonEgldOnlyEl = document.querySelector('button[name="toggle-egld-only"]')!;
                const mexInputsEl = document.querySelectorAll<HTMLInputElement>('#section-mex input')!;
                const egldPctEl = document.querySelector('input[name="egld-pct"]')! as HTMLInputElement;


                if (this._state.enabledMexAndEgld) {
                    sectionEl.classList.remove('invisible');
                    buttonEgldOnlyEl.classList.replace('bg-blue-200', 'bg-yellow-200');
                    buttonEgldOnlyEl.firstElementChild!.classList.replace('translate-x-5', 'translate-x-0');

                    egldPctEl.disabled = false;

                    mexInputsEl.forEach(el => {
                        el.required = true;
                    });
                } else {
                    sectionEl.classList.add('invisible');
                    buttonEgldOnlyEl.classList.replace('bg-yellow-200', 'bg-blue-200');
                    buttonEgldOnlyEl.firstElementChild!.classList.replace('translate-x-0', 'translate-x-5');

                    egldPctEl.value = '100';
                    egldPctEl.disabled = true;
                    egldPctEl.dispatchEvent(new Event('change'));


                    mexInputsEl.forEach(el => {
                        el.required = false;
                    });
                }
                break;
            case 'toggle-mex-locked-rewards':
                this._state.enabledMexRewardsLocked = !this._state.enabledMexRewardsLocked;
                const buttonLockedMexEl = document.querySelector('button[name="toggle-mex-locked-rewards"]')!;
                if (!this._state.enabledMexRewardsLocked) {
                    buttonLockedMexEl.classList.replace('bg-blue-200', 'bg-yellow-200');
                    buttonLockedMexEl.firstElementChild!.classList.replace('translate-x-5', 'translate-x-0');
                } else {
                    buttonLockedMexEl.classList.replace('bg-yellow-200', 'bg-blue-200');
                    buttonLockedMexEl.firstElementChild!.classList.replace('translate-x-0', 'translate-x-5');
                }
                break;
        }
    }

    handleSubmit() {
        console.debug('submit stake');
        const request = {} as { [key: string]: any };

        const egldInputsEl = document.querySelectorAll<HTMLInputElement>('#section-egld input')!;
        egldInputsEl.forEach(el => {
            let value = el.value;
            if (el.type === 'date') {
                const now = new Date();
                const target = new Date(value);

                request[el.name] = Math.floor((target.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
                return;
            }

            request[el.name] = value;
        });

        const egldSelectsEl = document.querySelectorAll<HTMLSelectElement>('#section-egld select')!;
        egldSelectsEl.forEach(el => {
            if (el.name === 'redelegation-interval') {
                request[el.name] = parseInt(el.value);
                return;
            }

            request[el.name] = el.value;
        });

        if (this._state.enabledMexAndEgld) {
            const mexInputsEl = document.querySelectorAll<HTMLInputElement>('#section-mex input')!;
            mexInputsEl.forEach(el => {
                request[el.name] = el.value;
            });

            request['mex-rewards-locked'] = this._state.enabledMexRewardsLocked;
        }

        console.debug("request payload", request);
        this.fetchSubmit(request);
    }

    connectedCallback() {
        console.debug('connected', this._state);

        this.fetchPrices();
        this.fetchStakingProviders();

        const formEl = document.querySelector('form')! as HTMLFormElement;
        formEl.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });

        const targetDateEl = document.querySelector('input[name="target-date-days"]')! as HTMLInputElement;
        targetDateEl.min = new Date().toISOString().split('T')[0];

        const selectsEl = document.querySelectorAll('select')!;
        selectsEl.forEach(selectEl => {
            selectEl.addEventListener('change', (e) => {
                const el = (e.target as HTMLSelectElement)
                this.handleInputEvent(el.name, el.value);
            });
        });

        const inputsEl = this.querySelectorAll('input')!;
        inputsEl.forEach(inputEl => {
            inputEl.addEventListener('change', (e) => {
                const el = (e.target as HTMLInputElement);
                this.handleInputEvent(el.name, el.value);
            });
        });

        const buttonsEl = this.querySelectorAll('button')!;
        buttonsEl.forEach(buttonEl => {
            buttonEl.addEventListener('click', (e) => {
                const el = (e.target as HTMLButtonElement);
                this.handleInputEvent(el.name, 'click');
            });
        });
    }
}
customElements.define('app-logic', AppLogic);


export class RewardDisplay extends HTMLElement {
    result: StrategyResult;
    strategyName: string;

    constructor(strategyName?: string, result?: StrategyResult) {
        super();
        this.result = result || {
            ProfitInEgld: "",
            ProfitInMex: "",
            ProfitInUSD: "",
            ROI: "",
            TotalBalanceInEgld: "",
            TotalBalanceInMex: "",
            TotalBalanceInUsd: "",
        };
        this.strategyName = strategyName || '';
    }

    connectedCallback() {
        const result = this.result;
        this.innerHTML = `
<div class="w-full">
    <div class="shadow-md p-0 rounded-lg">
        <p class="h-10 p-2 text-xs truncate bg-gray-500 rounded-t-lg"> ${this.strategyName} </p>
        <div class="pb-4 bg-gray-200 text-gray-900">
            <p class="py-2 truncate px-4"> ${result.ProfitInEgld} </p>
            <p class="py-2 truncate px-4 bg-white"> ${result.ProfitInMex} </p>
            <p class="py-2 truncate px-4"> ${result.ProfitInUSD} </p>
            <p class="py-2 truncate px-4 bg-white"> ${result.TotalBalanceInEgld} </p>
            <p class="py-2 truncate px-4"> ${result.TotalBalanceInMex} </p>
            <p class="py-2 truncate px-4 bg-white"> ${result.TotalBalanceInUsd} </p>
        </div>
        <div class="pt-2 pb-2 bg-gray-500 rounded-b-lg">
            <p class="font-extrabold text-2xl text-center"> +${result.ROI}% </p>
        </div>
    </div>
</div>`
    }
}

customElements.define('reward-display', RewardDisplay)
