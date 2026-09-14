#!/usr/bin/env node
/*
 * seed-personal-finance.mjs — seed realistic PERSONAL FINANCE test data.
 *
 * Populates the admin's personal "My Receipts" group with ~12 months of lifelike
 * transactions so the app can be used/tested as a personal finance tool: monthly
 * income, recurring bills (home loan, water, power, internet, phone, insurance),
 * groceries (with line items), dining, transport, subscriptions, shopping, health,
 * and one-off spends. Categorised + tagged so the Report Builder and pie chart show
 * meaningful breakdowns.
 *
 * Income is recorded as NEGATIVE amounts (a credit / money-in, same convention the
 * app uses for refunds), so a report grouped by category shows income as a negative
 * bucket and the grand total reads as NET cash flow (spend minus income).
 *
 * Auth: logs in as admin/admin (override with ADMIN_USERNAME/ADMIN_PASSWORD) and
 * uses the returned JWT as a Bearer token. No API key needed.
 *
 * Usage:  node api/dev/seed-personal-finance.mjs
 *
 * Env (all optional):
 *   API_BASE_URL     default http://localhost:8081/api
 *   ADMIN_USERNAME   default admin
 *   ADMIN_PASSWORD   default admin
 *   MONTHS           default 12      (how many trailing months to seed, ending this month)
 *   RANDOM_SEED      default 7
 *   GROUP_NAME       default "My Receipts"
 *
 * Requires Node 18+ (built-in fetch). No npm install.
 */

const API_BASE_URL = (process.env.API_BASE_URL || "http://localhost:8081/api").replace(/\/+$/, "");
const ADMIN_USERNAME = process.env.ADMIN_USERNAME || "admin";
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || "admin";
const MONTHS = Number.parseInt(process.env.MONTHS ?? "12", 10) || 12;
const GROUP_NAME = process.env.GROUP_NAME || "My Receipts";

// ---------- seeded PRNG (mulberry32) ----------
let _rng = (Number.parseInt(process.env.RANDOM_SEED ?? "7", 10) || 7) >>> 0;
function rand() {
  _rng = (_rng + 0x6d2b79f5) | 0;
  let t = _rng;
  t = Math.imul(t ^ (t >>> 15), t | 1);
  t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
}
const randInt = (min, max) => min + Math.floor(rand() * (max - min + 1));
const pick = (arr) => arr[Math.floor(rand() * arr.length)];
const money = (min, max) => (min + rand() * (max - min)).toFixed(2);
const jitter = (base, pct) => (base * (1 + (rand() * 2 - 1) * pct)).toFixed(2);

// ---------- HTTP ----------
let JWT = "";
async function api(method, path, body) {
  const res = await fetch(API_BASE_URL + path, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(JWT ? { Authorization: `Bearer ${JWT}` } : {}),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const text = await res.text();
  if (res.status >= 400) {
    throw new Error(`${method} ${path} -> ${res.status}: ${text.slice(0, 240)}`);
  }
  return text ? JSON.parse(text) : null;
}

// A date on a given month offset back from now, on a specific day-of-month.
function dateFor(monthsAgo, day) {
  const now = new Date();
  const d = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth() - monthsAgo, day, 0, 0, 0));
  return d.toISOString().replace(/\.\d{3}Z$/, "Z");
}

// ---------- catalog ----------
// Expenses are positive; income is negative (a credit). "recurring" bills get one
// entry per month at a jittered amount; "variable" categories get a random count
// of entries per month within an amount range.
const CATEGORIES = [
  { name: "Income", description: "Salary and other money coming in" },
  { name: "Home Loan", description: "Mortgage repayments" },
  { name: "Rent", description: "Rent payments" },
  { name: "Water", description: "Water utility bill" },
  { name: "Electricity", description: "Power bill" },
  { name: "Gas", description: "Natural gas bill" },
  { name: "Internet", description: "Home broadband" },
  { name: "Phone", description: "Mobile phone plan" },
  { name: "Insurance", description: "Home, car, and health insurance" },
  { name: "Groceries", description: "Supermarket and food shopping" },
  { name: "Dining Out", description: "Restaurants, cafes, takeaway" },
  { name: "Transport", description: "Fuel, public transport, rideshare" },
  { name: "Subscriptions", description: "Streaming and software subscriptions" },
  { name: "Health", description: "Pharmacy, doctor, dentist" },
  { name: "Shopping", description: "Clothes, electronics, household" },
  { name: "Entertainment", description: "Movies, events, hobbies" },
  { name: "Savings", description: "Transfers to savings / investment" },
];

const TAGS = ["Essential", "Discretionary", "Recurring", "Income", "Tax Deductible"];

const GROCERY_STORES = ["Countdown", "New World", "PAK'nSAVE", "Asian Supermarket", "Costco"];
const GROCERY_ITEMS = [
  ["Milk 2L", 4.2], ["Bread", 3.5], ["Eggs (dozen)", 8.9], ["Chicken breast", 12.5],
  ["Rice 5kg", 15.9], ["Bananas", 3.2], ["Coffee beans", 18.5], ["Cheese block", 11.0],
  ["Pasta", 2.8], ["Tomatoes", 5.4], ["Yoghurt", 6.2], ["Dish soap", 4.9],
];
const RESTAURANTS = ["Local Cafe", "Thai Kitchen", "Sushi Bar", "Pizza Place", "Burger Joint", "Ramen House"];
const FUEL = ["BP", "Z Energy", "Mobil", "Gull"];
const SUBS = [["Netflix", 24.99], ["Spotify", 14.99], ["Disney+", 12.99], ["iCloud", 4.99], ["YouTube Premium", 15.99]];
const SHOPS = ["Kmart", "The Warehouse", "JB Hi-Fi", "Farmers", "Mitre 10", "Briscoes"];

async function main() {
  console.log(`Seeding personal-finance data against ${API_BASE_URL}`);

  const login = await api("POST", "/login/?tokensInBody=true", {
    username: ADMIN_USERNAME,
    password: ADMIN_PASSWORD,
  });
  JWT = login.jwt;
  if (!JWT) throw new Error("login did not return a jwt");

  const appData = await api("GET", "/user/appData");
  const adminId = appData.claims.userId;
  const group = (appData.groups || []).find((g) => g.name === GROUP_NAME && !g.isAllGroup);
  if (!group) throw new Error(`could not find personal group "${GROUP_NAME}"`);
  console.log(`Authenticated as ${appData.claims.username} (id ${adminId}); group "${GROUP_NAME}" (id ${group.id})`);

  // Categories (idempotent by name).
  const catByName = new Map((appData.categories || []).map((c) => [c.name, c.id]));
  for (const cat of CATEGORIES) {
    if (!catByName.has(cat.name)) {
      const created = await api("POST", "/category", { name: cat.name, description: cat.description });
      catByName.set(cat.name, created.id);
      console.log(`  + category ${cat.name}`);
    }
  }
  const cat = (name) => ({ id: catByName.get(name), name });

  // Tags (idempotent by name).
  const tagByName = new Map((appData.tags || []).map((t) => [t.name, t.id]));
  for (const name of TAGS) {
    if (!tagByName.has(name)) {
      const created = await api("POST", "/tag", { name, description: "" });
      tagByName.set(name, created.id);
      console.log(`  + tag ${name}`);
    }
  }
  const tag = (name) => ({ id: tagByName.get(name), name });

  const bodies = [];
  const base = { groupId: group.id, paidByUserId: adminId, status: "RESOLVED" };

  for (let m = MONTHS - 1; m >= 0; m--) {
    // Income: salary mid- and end-of-month (negative = money in).
    bodies.push({ ...base, name: "Salary", date: dateFor(m, 15), amount: `-${jitter(4300, 0.02)}`,
      categories: [cat("Income")], tags: [tag("Income")] });
    bodies.push({ ...base, name: "Salary", date: dateFor(m, 30 > 28 ? 28 : 30), amount: `-${jitter(4300, 0.02)}`,
      categories: [cat("Income")], tags: [tag("Income")] });

    // Recurring bills.
    bodies.push({ ...base, name: "Home Loan Repayment", date: dateFor(m, 2), amount: jitter(2380, 0.0),
      categories: [cat("Home Loan")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Watercare", date: dateFor(m, 8), amount: jitter(64, 0.15),
      categories: [cat("Water")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Electricity Bill", date: dateFor(m, 10), amount: jitter(150, 0.35),
      categories: [cat("Electricity")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Gas Bill", date: dateFor(m, 10), amount: jitter(55, 0.3),
      categories: [cat("Gas")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Broadband", date: dateFor(m, 5), amount: "89.00",
      categories: [cat("Internet")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Mobile Plan", date: dateFor(m, 5), amount: "45.00",
      categories: [cat("Phone")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Insurance Premium", date: dateFor(m, 12), amount: jitter(185, 0.05),
      categories: [cat("Insurance")], tags: [tag("Essential"), tag("Recurring")] });
    bodies.push({ ...base, name: "Savings Transfer", date: dateFor(m, 16), amount: "500.00",
      categories: [cat("Savings")], tags: [tag("Recurring")] });

    // Subscriptions (a couple each month).
    for (const [svc, price] of SUBS.filter(() => rand() < 0.6)) {
      bodies.push({ ...base, name: svc, date: dateFor(m, randInt(3, 20)), amount: price.toFixed(2),
        categories: [cat("Subscriptions")], tags: [tag("Discretionary"), tag("Recurring")] });
    }

    // Groceries: 3–5 per month, some with line items.
    const groceryTrips = randInt(3, 5);
    for (let i = 0; i < groceryTrips; i++) {
      const store = pick(GROCERY_STORES);
      const itemCount = randInt(3, 6);
      const items = [];
      let total = 0;
      for (let j = 0; j < itemCount; j++) {
        const [nm, unit] = pick(GROCERY_ITEMS);
        const qty = randInt(1, 3);
        const amount = +(unit * qty).toFixed(2);
        total += amount;
        items.push({
          name: nm, quantity: qty.toString(), unitPrice: unit.toFixed(2),
          amount: amount.toFixed(2), status: "OPEN",
        });
      }
      bodies.push({ ...base, name: store, date: dateFor(m, randInt(1, 27)), amount: total.toFixed(2),
        categories: [cat("Groceries")], tags: [tag("Essential")], receiptItems: items });
    }

    // Dining out: 2–5 per month.
    for (let i = 0; i < randInt(2, 5); i++) {
      bodies.push({ ...base, name: pick(RESTAURANTS), date: dateFor(m, randInt(1, 27)),
        amount: money(14, 95), categories: [cat("Dining Out")], tags: [tag("Discretionary")] });
    }

    // Transport / fuel: 2–4 per month.
    for (let i = 0; i < randInt(2, 4); i++) {
      bodies.push({ ...base, name: pick(FUEL), date: dateFor(m, randInt(1, 27)),
        amount: money(45, 110), categories: [cat("Transport")], tags: [tag("Essential")] });
    }

    // Occasional health, shopping, entertainment.
    if (rand() < 0.5) bodies.push({ ...base, name: "Pharmacy", date: dateFor(m, randInt(1, 27)),
      amount: money(12, 90), categories: [cat("Health")], tags: [tag("Essential")] });
    if (rand() < 0.7) bodies.push({ ...base, name: pick(SHOPS), date: dateFor(m, randInt(1, 27)),
      amount: money(25, 260), categories: [cat("Shopping")], tags: [tag("Discretionary")] });
    if (rand() < 0.6) bodies.push({ ...base, name: pick(["Movie Night", "Concert", "Bowling", "Mini Golf"]),
      date: dateFor(m, randInt(1, 27)), amount: money(20, 140),
      categories: [cat("Entertainment")], tags: [tag("Discretionary")] });
  }

  console.log(`Creating ${bodies.length} transactions over ${MONTHS} months...`);
  let ok = 0, failed = 0;
  for (const body of bodies) {
    try { await api("POST", "/receipt", body); ok++; }
    catch (err) { failed++; if (failed <= 5) console.error(`  failed: ${err.message}`); }
    if ((ok + failed) % 50 === 0) console.log(`  ${ok + failed}/${bodies.length}`);
  }

  console.log("\nDone.");
  console.log(`  group:        ${GROUP_NAME} (id ${group.id})`);
  console.log(`  categories:   ${CATEGORIES.length}, tags: ${TAGS.length}`);
  console.log(`  transactions: ok ${ok}, failed ${failed}`);
  console.log(`  income is negative (money in); expenses positive; grand total = net cash flow.`);
}

main().catch((err) => { console.error("FATAL:", err.message); process.exit(1); });
