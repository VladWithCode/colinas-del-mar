const PRICE_NO_DEPOSIT = 245_000;
const PRICE_WITH_DEPOSIT = 230_000;
const DEC_PAYMENT_NO_DEPOSIT = 11_250;
const DEC_PAYMENT_WITH_DEPOSIT = 10_000;
const DEPOSIT = 30_000;

const formatter = new Intl.NumberFormat("es-MX", { style: "currency", currency: "MXN" });

const depositOptionInp = document.querySelectorAll("[data-radio-deposit]");
const quoteNoPaymentsEl = document.getElementById("quote-no-payments-value");
const quoteNoPaymentsRange = document.getElementById("quote-no-payments-range");
const quoteMonthlyEl = document.getElementById("quote-monthly");
const quoteDecPaymentEl = document.getElementById("quote-december");

let withDeposit = true;
let selectedPrice = PRICE_WITH_DEPOSIT;
let selectedDecPayment = DEC_PAYMENT_WITH_DEPOSIT;
let nPayments = 44;
let monthlyPayment = calculateDecemberPayment(withDeposit, nPayments);

for (const depositOption of depositOptionInp) {
    depositOption.addEventListener("change", function(e) {
        withDeposit = e.target.value === "with-deposit";

        if (withDeposit) {
            selectedPrice = PRICE_WITH_DEPOSIT;
            selectedDecPayment = DEC_PAYMENT_WITH_DEPOSIT;
        } else {
            selectedPrice = PRICE_NO_DEPOSIT;
            selectedDecPayment = DEC_PAYMENT_NO_DEPOSIT;
        }

        quoteMonthlyEl.textContent = calculateDecemberPayment();
        quoteDecPaymentEl.textContent = formatter.format(selectedDecPayment);
    });
}

quoteNoPaymentsRange.addEventListener("input", function() {
    quoteNoPaymentsEl.innerHTML = this.value;
    nPayments = this.value;
    quoteMonthlyEl.textContent = calculateDecemberPayment();
});

function calculateDecemberPayment() {
    let amountPending = 0

    if (withDeposit) {
        amountPending = PRICE_WITH_DEPOSIT - DEPOSIT - (DEC_PAYMENT_WITH_DEPOSIT * 4);
    } else {
        amountPending = PRICE_NO_DEPOSIT - (DEC_PAYMENT_NO_DEPOSIT * 4);
    }

    return formatter.format(amountPending / nPayments);
}
