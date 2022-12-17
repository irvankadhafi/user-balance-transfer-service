# user-balance-transfer-service

In the context of banking, the term "debit" typically refers to a transaction that reduces the balance of an account or increases debt. On the other hand, the term "credit" typically refers to a transaction that increases the balance of an account or reduces debt. Therefore, it can be said that a debit is a withdrawal of money while a credit is a receipt of money.
<div class="markdown prose break-words dark:prose-invert dark">

## To-Do

*   Add unit tests
*   Implement transfer between banks
*   Implement transfer from user to bank and vice versa
*   Implement refresh token
*   Implement limit of only 1 session per user
*   Improve clean code

# API Documentation

Below are the available API endpoints for the project:

## Migrations and Seeds

*   To run migrations, execute `make migrate`.
*   To initialize user data, run `make seed-user`.
*   To start the server, run `make run`.

## Authentication

To log in, send a `POST` request to `localhost:3000/auth/login` with the following body:

<pre>{
    <span class="hljs-string">"email"</span>: <span class="hljs-string">"johndoe@mail.com"</span>,
    <span class="hljs-string">"password"</span>: <span class="hljs-string">"123456"</span>
}
</pre>

## User Balance

### Add balance

To add balance to a user's account, send a `POST` request to `localhost:3000/user-balance/add` with the following JSON body:

<pre>
{
    <span class="hljs-attr">"balance"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"10000"</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"author"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"irvan"</span>
<span class="hljs-punctuation">}
</pre>

If successful, the API will return an `OK` response.

### Get balance details

To get the details of a user's balance, send a `GET` request to `localhost:3000/user-balance/`. The API will return a JSON object with the following structure:

<pre>
{</span>
    <span class="hljs-attr">"id"</span><span class="hljs-punctuation">:</span> <span class="hljs-number">80</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"user_id"</span><span class="hljs-punctuation">:</span> <span class="hljs-number">1</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"balance"</span><span class="hljs-punctuation">:</span> <span class="hljs-number">10000</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"balance_achieve"</span><span class="hljs-punctuation">:</span> <span class="hljs-number">40000</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"created_at"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"2022-12-17T07:26:28.27784Z"</span>
<span class="hljs-punctuation">}
</pre>

### Transfer balance between users

To transfer balance between users, send a `POST` request to `localhost:3000/user-balance/transfer/` with the following JSON body:

<pre>
{
    <span class="hljs-attr">"to_user_id"</span> <span class="hljs-punctuation">:</span> <span class="hljs-number">2</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"balance"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"100000"</span><span class="hljs-punctuation">,</span>
    <span class="hljs-attr">"author"</span><span class="hljs-punctuation">:</span> <span class="hljs-string">"irvan"</span>
<span class="hljs-punctuation">}
</pre>

If successful, the API will return an `OK` response.

## Bank Balance

### Create bank account

To create a bank account, send a `POST` request to `localhost:3000/bank-balance/create` with the following JSON body:

<pre>
{
    "<span class="hljs-selector-tag">code</span>":<span class="hljs-string">"ABCDE12345"</span>
}` 
</pre>

If successful, the API will return an `OK` response.

### Add balance to bank account

To add balance to a bank account, send a `POST` request to `localhost:3000/bank-balance/add` with the following JSON body:

<pre>
{
    "<span class="hljs-selector-tag">code</span>": <span class="hljs-string">"ABCDE12345"</span>,
    <span class="hljs-string">"balance"</span>: <span class="hljs-string">"10000"</span>,
    <span class="hljs-string">"author"</span>: <span class="hljs-string">"irvan"</span>
}
</pre>

If successful, the API will return an `OK` response.



</div>